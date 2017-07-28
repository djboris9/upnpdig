package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/djboris9/go-upnp/description"
	"github.com/djboris9/go-upnp/schemas"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe UPnP devices and services on your network",
	Long: `Describe is used to show properties of your UPnP devices and show
description about the services like actions`,

	Run: func(cmd *cobra.Command, args []string) {
		// Parse flags
		deviceUrl, err := cmd.Flags().GetString("device")
		if err != nil {
			log.Fatal(err)
		}
		serviceId, err := cmd.Flags().GetString("service")
		if err != nil {
			log.Fatal(err)
		}

		device, err := description.GetDeviceFromURL(deviceUrl)
		if err != nil {
			log.Fatal(err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		if serviceId == "" {
			// Print device
			describeDevicePrinter(w, device.Device, "")
		} else {
			// Print service of device according to ServiceId
			serviceLoc, err := describeGetServiceLocation(deviceUrl, device.Device, serviceId)
			if err != nil {
				log.Fatal(err)
			}

			service, err := description.GetServiceFromURL(serviceLoc.String())
			if err != nil {
				log.Fatal(err)
			}

			describeServicePrinter(w, service)
		}
		w.Flush()
	},
}

var describeErrServiceNotFound = errors.New("Service not found")

func describeGetServiceLocation(deviceUrl string, dev schemas.DeviceType, serviceId string) (*url.URL, error) {
	for _, svc := range dev.ServiceList.Services {
		if svc.ServiceId == serviceId {
			// Resolve service URL
			devUrl, err := url.Parse(deviceUrl)
			if err != nil {
				return nil, err
			}

			svcUrl, err := url.Parse(svc.SCPDURL)
			if err != nil {
				return nil, err
			}

			return devUrl.ResolveReference(svcUrl), nil
		}
	}

	for _, embeddedDevice := range dev.DeviceList.Devices {
		u, err := describeGetServiceLocation(deviceUrl, embeddedDevice, serviceId)
		if err == nil {
			return u, nil
		} else if err != describeErrServiceNotFound {
			log.Printf("%v", err)
		}
	}

	return nil, describeErrServiceNotFound
}

func describeServicePrinter(w *tabwriter.Writer, svc schemas.Service) {
	fmt.Fprintln(w, "State variables")
	for _, sv := range svc.ServiceStateTable.StateVariable {
		fmt.Fprintf(w, "Name\t%v\n", sv.Name)
		fmt.Fprintf(w, "Multicast\t%v\n", sv.Multicast)
		fmt.Fprintf(w, "SendEvents\t%v\n", sv.SendEvents)
		fmt.Fprintf(w, "DefaultValue\t%v\n", sv.DefaultValue)
		fmt.Fprintf(w, "DataType\t%v\n", sv.DataType)
		fmt.Fprintf(w, "AllowedValueList\t%v\n", sv.AllowedValueList)
		fmt.Fprintf(w, "AllowedValueRange\t%v-%v, step %v\n", sv.AllowedValueRange.Minimum, sv.AllowedValueRange.Maximum, sv.AllowedValueRange.Step)
		fmt.Fprintln(w, "")
	}

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Actions")
	for _, action := range svc.ActionList.Action {
		fmt.Fprintf(w, "Name\t%v\n", action.Name)
		fmt.Fprintln(w, "Arguments:")
		for _, arg := range action.ArgumentList.Argument {
			fmt.Fprintf(w, "\tName\t%v\n", arg.Name)
			fmt.Fprintf(w, "\tDirection\t%v\n", arg.Direction)
			fmt.Fprintf(w, "\tRelatedStateVariable\t%v\n", arg.RelatedStateVariable)
			fmt.Fprintf(w, "\tRetval\t%v\n", arg.Retval)
			fmt.Fprintln(w, "")
		}
		fmt.Fprintln(w, "")
	}
}

func describeDevicePrinter(w *tabwriter.Writer, device schemas.DeviceType, padding string) {
	fmt.Fprintf(w, "%sDeviceType\t%v\n", padding, device.DeviceType)
	fmt.Fprintf(w, "%sFriendlyName\t%v\n", padding, device.FriendlyName)
	fmt.Fprintf(w, "%sManufacturer\t%v\n", padding, device.Manufacturer)
	fmt.Fprintf(w, "%sManufacturerURL\t%v\n", padding, device.ManufacturerURL)
	fmt.Fprintf(w, "%sModelName\t%v\n", padding, device.ModelName)
	fmt.Fprintf(w, "%sModelDescription\t%v\n", padding, device.ModelDescription)
	fmt.Fprintf(w, "%sModelURL\t%v\n", padding, device.ModelURL)
	fmt.Fprintf(w, "%sSerialNumber\t%v\n", padding, device.SerialNumber)
	fmt.Fprintf(w, "%sPresentationURL\t%v\n", padding, device.PresentationURL)
	fmt.Fprintf(w, "%sUPC\t%v\n", padding, device.UPC)
	fmt.Fprintf(w, "%sUDN\t%v\n", padding, device.UDN)
	fmt.Fprintln(w, "")

	if len(device.ServiceList.Services) > 0 {
		fmt.Fprintf(w, "%sServices\n", padding)
		for _, svc := range device.ServiceList.Services {
			fmt.Fprintf(w, "\t%sServiceId\t%v\n", padding, svc.ServiceId)
			fmt.Fprintf(w, "\t%sServiceType\t%v\n", padding, svc.ServiceType)
			fmt.Fprintf(w, "\t%sControlURL\t%v\n", padding, svc.ControlURL)
			fmt.Fprintf(w, "\t%sEventSubURL\t%v\n", padding, svc.EventSubURL)
			fmt.Fprintf(w, "\t%sSCPD URL\t%v\n", padding, svc.SCPDURL)
			fmt.Fprintln(w, "")
		}
	}

	if len(device.DeviceList.Devices) > 0 {
		fmt.Fprintf(w, "%sDevices:\n", padding)
		for _, dev2 := range device.DeviceList.Devices {
			describeDevicePrinter(w, dev2, padding+"    ")
			fmt.Fprintln(w, "")
		}
	}
}

func init() {
	RootCmd.AddCommand(describeCmd)
	describeCmd.PersistentFlags().StringP("device", "d", "http://localhost:1400/device_description.xml", "Device description URL")
	describeCmd.PersistentFlags().StringP("service", "s", "", "ServiceId to describe")
}
