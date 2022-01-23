#!/bin/bash

set -euxo pipefail

RTL_SDR_DEV_PATH=${RTL_SDR_DEV_PATH:=/dev/rtl_sdr}

FREQUENCY=${FREQUENCY:=433M}

RTL_433_FLAGS=${RTL_433_FLAGS:=}

MATCHERS=${MATCHERS:=}

SUBPROCESS="rtl_433 -F json -f $FREQUENCY $RTL_433_FLAGS"

if [ -e "$RTL_SDR_DEV_PATH" ]; then
    # Here we remap /dev/rtl_sdr to the correct path in /dev/bus/usb so rtl_433 can discover it correctly

    # Use udevadm to discover the BUSNUM and DEVNUM
    BUSNUM=$(udevadm info --name /dev/rtl_sdr | grep BUSNUM | awk -F = '{ print $2 }')
    DEVNUM=$(udevadm info --name /dev/rtl_sdr | grep DEVNUM | awk -F = '{ print $2 }')

    DEVDIR=/dev/bus/usb/$BUSNUM
    DEVPATH=$DEVDIR/$DEVNUM

    # Create the correct path to the eventual device
    mkdir -p /dev/bus/usb/$BUSNUM

    # Symlink original device to correct path
    ln -s $RTL_SDR_DEV_PATH $DEVPATH
fi

/rtl_433_prometheus --subprocess "$SUBPROCESS" $MATCHERS
