package main

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestChannel(t *testing.T) {
	cases := []struct {
		in   Message
		want string
	}{
		{Message{RawChannel: ""}, ""},
		{Message{RawChannel: "1"}, "1"},
		{Message{RawChannel: "2"}, "2"},
		{Message{RawChannel: 0.0}, "0"},
		{Message{RawChannel: 2.0}, "2"},
	}
	for _, tt := range cases {
		msg := tt.in
		want := tt.want
		got, err := msg.Channel()
		if err != nil {
			t.Errorf("unexpected err=%v", err)
		}
		if got != want {
			t.Errorf("%+v.Channel()=%v, want=%v", msg, got, want)
		}
	}
}

func TestParsingToMetrics(t *testing.T) {
	fn := "test_input.txt"
	f, err := os.Open(fn)
	if err != nil {
		t.Fatalf("couldn't open %v: %v", fn, err)
	}
	defer f.Close()
	err = run(f)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	wantTemperature := `
		# HELP rtl_433_temperature_celsius Temperature in Celsius
		# TYPE rtl_433_temperature_celsius gauge
 		rtl_433_temperature_celsius{channel="1",id="77",location="",model="AmbientWeather-TX8300"} 16.5
		rtl_433_temperature_celsius{channel="1",id="94",location="",model="Nexus-TH"} 22.6
		rtl_433_temperature_celsius{channel="2",id="184",location="",model="Nexus-TH"} 21.7
		rtl_433_temperature_celsius{channel="2",id="59",location="",model="Ambientweather-F007TH"} 33.166666666666664
		rtl_433_temperature_celsius{channel="3",id="55",location="",model="Ecowitt-WH53"} 18
		rtl_433_temperature_celsius{channel="A",id="7997",location="",model="Acurite-Tower"} 12.6
		rtl_433_temperature_celsius{channel="",id="52451",location="",model="Holman-WS5029"} 26.200
	`

	if err := testutil.CollectAndCompare(temperature, strings.NewReader(wantTemperature), "rtl_433_temperature_celsius"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}

	wantHumidity := `
		# HELP rtl_433_humidity Relative Humidity (0-1.0)
		# TYPE rtl_433_humidity gauge
		rtl_433_humidity{channel="1",id="77",location="",model="AmbientWeather-TX8300"} 0.66
		rtl_433_humidity{channel="1",id="94",location="",model="Nexus-TH"} 0.53
		rtl_433_humidity{channel="2",id="184",location="",model="Nexus-TH"} 0.55
		rtl_433_humidity{channel="2",id="59",location="",model="Ambientweather-F007TH"} 0.35
		rtl_433_humidity{channel="A",id="7997",location="",model="Acurite-Tower"} 0.91
		rtl_433_humidity{channel="",id="52451",location="",model="Holman-WS5029"} 0.59
	`
	if err := testutil.CollectAndCompare(humidity, strings.NewReader(wantHumidity), "rtl_433_humidity"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}

	wantPacketsReceived := `
		# HELP rtl_433_packets_received Packets (temperature messages) received.
		# TYPE rtl_433_packets_received counter
		rtl_433_packets_received{channel="1",id="77",location="",model="AmbientWeather-TX8300"} 1
		rtl_433_packets_received{channel="1",id="94",location="",model="Nexus-TH"} 1
		rtl_433_packets_received{channel="2",id="184",location="",model="Nexus-TH"} 1
		rtl_433_packets_received{channel="2",id="59",location="",model="Ambientweather-F007TH"} 1
		rtl_433_packets_received{channel="3",id="55",location="",model="Ecowitt-WH53"} 1
		rtl_433_packets_received{channel="A",id="7997",location="",model="Acurite-Tower"} 1
		rtl_433_packets_received{channel="",id="52451",location="",model="Holman-WS5029"} 1
	`
	if err := testutil.CollectAndCompare(packetsReceived, strings.NewReader(wantPacketsReceived), "rtl_433_packets_received"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
	wantBattery := `
		# HELP rtl_433_battery Battery high (1) or low (0).
		# TYPE rtl_433_battery gauge
		rtl_433_battery{channel="1",id="94",location="",model="Nexus-TH"} 1
		rtl_433_battery{channel="2",id="184",location="",model="Nexus-TH"} 1
		rtl_433_battery{channel="2",id="59",location="",model="Ambientweather-F007TH"} 1
		rtl_433_battery{channel="A",id="7997",location="",model="Acurite-Tower"} 0
	`
	if err := testutil.CollectAndCompare(battery, strings.NewReader(wantBattery), "rtl_433_battery"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}

	wantRain := `
		# HELP rtl_433_rain_mm Rain in mm
		# TYPE rtl_433_rain_mm gauge
		rtl_433_rain_mm{channel="",id="52451",location="",model="Holman-WS5029"} 28.440
	`
	if err := testutil.CollectAndCompare(rain, strings.NewReader(wantRain), "rtl_433_rain_mm"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}

	wantWindSpeed := `
		# HELP rtl_433_wind_avg_km_h Average wind speed in km/h
		# TYPE rtl_433_wind_avg_km_h gauge
		rtl_433_wind_avg_km_h{channel="",id="52451",location="",model="Holman-WS5029"} 11
	`
	if err := testutil.CollectAndCompare(wind_speed, strings.NewReader(wantWindSpeed), "rtl_433_wind_avg_km_h"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}

	wantWindDirection := `
		# HELP rtl_433_wind_dir_deg Wind direction in degrees
		# TYPE rtl_433_wind_dir_deg gauge
		rtl_433_wind_dir_deg{channel="",id="52451",location="",model="Holman-WS5029"} 135
	`
	if err := testutil.CollectAndCompare(wind_direction, strings.NewReader(wantWindDirection), "rtl_433_wind_dir_deg"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}
