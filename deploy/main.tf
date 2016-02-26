provider "digitalocean" {}

variable "ssh_keys" {}

variable "region"       { default="nyc3" }
variable "web_size"     { default="4gb" }
variable "metrics_size" { default="1gb" }
variable "brog_version" { default="0.3" }
variable "prom_version" { default="0.17.0rc2" }

resource "digitalocean_droplet" "web" {
    count = 1

    image              = "debian-8-x64"
    name               = "web${count.index}.aybabt.me"
    region             = "${var.region}"
    size               = "${var.web_size}"
    ssh_keys           = ["${split(",", var.ssh_keys)}"]
    private_networking = true

    provisioner "remote-exec" {
            inline = <<CMD
apt-get install -y git

mkdir /srv/
git clone https://github.com/aybabtme/antoine.im.git /srv/aybabt.me

mkdir -p /opt/bin/brog/
wget -qO- https://github.com/aybabtme/brog/releases/download/${var.brog_version}/brog_linux.tar.gz | tar xvz
mv ./brog /opt/bin/brog/brog

cat > /etc/systemd/system/web.service <<EOF
[Unit]
Description=aybabt.me
After=network.target

[Service]
Type=simple
WorkingDirectory=/srv/aybabt.me/
ExecStart=/opt/bin/brog/brog server prod
Restart=always

[Install]
WantedBy=multi-user.target
EOF

systemctl enable web.service
systemctl start web.service
CMD
    }
}


resource "digitalocean_droplet" "prometheus" {
    image              = "debian-8-x64"
    name               = "prometheus.aybabt.me"
    region             = "${var.region}"
    size               = "${var.metrics_size}"
    ssh_keys           = ["${split(",", var.ssh_keys)}"]
    private_networking = true

    provisioner "remote-exec" {
            inline = <<CMD
mkdir -p /opt/bin/prometheus/
wget -qO- https://github.com/prometheus/prometheus/releases/download/${var.prom_version}/prometheus-${var.prom_version}.linux-amd64.tar.gz | tar xvz
cd prometheus-*
mv ./prometheus /opt/bin/prometheus/prometheus

mkdir -p /var/lib/prometheus
mkdir -p /etc/prometheus

cat > /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=prometheus
After=network.target

[Service]
Type=simple
ExecStart=/opt/bin/prometheus/prometheus \
  -config.file=/etc/prometheus/config.yml \
  -storage.local.path=/var/lib/prometheus \
  -web.listen-address "localhost:9090" \
  -storage.local.retention 8640h \
  -log.format logger:stderr?json=true
Restart=always

[Install]
WantedBy=multi-user.target
EOF

cat > /etc/prometheus/config.yml <<EOF
global:
  scrape_interval:     2s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'aybabtme'

    scrape_interval: 2s
    scrape_timeout: 10s
    metrics_path: /debug/metrics

    target_groups:
      - targets: ['${digitalocean_droplet.web.ipv4_address_private}:80']

  - job_name: 'prometheus'

    scrape_interval: 5s
    scrape_timeout: 10s

    target_groups:
      - targets: ['localhost:9090']
EOF

systemctl enable prometheus.service
systemctl start prometheus.service
CMD
    }
}
