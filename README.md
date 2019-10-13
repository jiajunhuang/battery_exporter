# Battery Exporter

Battery Exporter for Prometheus, you can use this exporter to monitor your battery infomation include:

- energy by design
- energy now
- energy full
- battery charge cycle count

just like this:

```bash
$ http :9119/metrics | grep battery
# HELP battery_cycle_count Battery Charge Cycle Count
# TYPE battery_cycle_count gauge
battery_cycle_count 22
# HELP battery_energy_full Energy Full in mWh
# TYPE battery_energy_full gauge
battery_energy_full 5.665e+07
# HELP battery_energy_full_design Energy Full in mWh By Design
# TYPE battery_energy_full_design gauge
battery_energy_full_design 5.6316e+07
# HELP battery_energy_now Energy Now in mWh
# TYPE battery_energy_now gauge
battery_energy_now 4.3806e+07
```

## Usage

- download or compile binary yourself, and put it on `/usr/local/bin/battery_exporter`
- change the owner and permission: `sudo chown nobody:nobody /usr/local/bin/battery_exporter && sudo chmod +x /usr/local/bin/battery_exporter`
- copy the systemd service file to `/etc/systemd/system/prometheus-battery-exporter.service`
- start and enable it `sudo systemctl start prometheus-battery-exporter.service && sudo systemctl enable prometheus-battery-exporter.service`
