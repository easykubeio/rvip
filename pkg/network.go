package pkg

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type NetworkConfigurator interface {
	AddIP() error
	DeleteIP() error
	IsSet() (bool, error)
	IsUsed() bool
	IP() string
	Interface() string
	SendARP() error
}

type NetlinkNetworkConfigurator struct {
	address *netlink.Addr
	link    netlink.Link
}

func NewNetlinkNetworkConfigurator(_address, _interface string) (result NetlinkNetworkConfigurator, error error) {
	result = NetlinkNetworkConfigurator{}

	result.address, error = netlink.ParseAddr(_address + "/32")
	if error != nil {
		error = errors.Wrapf(error, "could not parse address '%s'", _address)

		return
	}

	result.link, error = netlink.LinkByName(_interface)
	if error != nil {
		error = errors.Wrapf(error, "could not get link for interface '%s'", _interface)

		return
	}

	return
}

func (configurator NetlinkNetworkConfigurator) AddIP() error {
	result, error := configurator.IsSet()
	if error != nil {
		return errors.Wrap(error, "ip check in AddIP failed")
	}

	// Already set
	if result {
		return nil
	}

	if error = netlink.AddrAdd(configurator.link, configurator.address); error != nil {
		return errors.Wrap(error, "could not add ip")
	}

	return nil
}

func (configurator NetlinkNetworkConfigurator) DeleteIP() error {
	result, error := configurator.IsSet()
	if error != nil {
		return errors.Wrap(error, "ip check in DeleteIP failed")
	}

	// Nothing to delete
	if !result {
		return nil
	}

	if error = netlink.AddrDel(configurator.link, configurator.address); error != nil {
		return errors.Wrap(error, "could not delete ip")
	}

	return nil
}

func (configurator NetlinkNetworkConfigurator) IsSet() (result bool, error error) {
	var addresses []netlink.Addr

	addresses, error = netlink.AddrList(configurator.link, 0)
	if error != nil {
		error = errors.Wrap(error, "could not list addresses")

		return
	}

	for _, address := range addresses {
		if address.Equal(*configurator.address) {
			return true, nil
		}
	}

	return false, nil
}

func (configurator NetlinkNetworkConfigurator) IP() string {
	return configurator.address.IP.String()
}

func (configurator NetlinkNetworkConfigurator) Interface() string {
	return configurator.link.Attrs().Name
}

func (configurator NetlinkNetworkConfigurator) IsUsed() bool {
	cmd := fmt.Sprintf("ping -w 1 -c 1 %s > /dev/null && echo true || echo false", configurator.IP())
	output, err := exec.Command("/bin/bash", "-c", cmd).Output()
	out := strings.TrimSpace(string(output))
	if err != nil || out == "false" {
		return false
	}
	return true
}

func (configurator NetlinkNetworkConfigurator) SendARP() error {
	var arp *ARPGratuitous
	arp = &ARPGratuitous{
		IfaceName: configurator.Interface(),
		IP: net.ParseIP(configurator.IP()),
	}
	if err := ARPSendGratuitous(arp); err != nil {
		return err
	}
	return nil
}
