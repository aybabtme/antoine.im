provider "digitalocean" {}

variable "region"   {}
variable "web_size" {}
variable "ssh_keys" {}
variable "brog_version" {default="0.3"}

resource "digitalocean_droplet" "web" {
    count = 1

    image              = "debian-8-x64"
    name               = "web${count.index}.aybabt.me.debian"
    region             = "${var.region}"
    size               = "${var.web_size}"
    ssh_keys           = ["${split(",", var.ssh_keys)}"]
    private_networking = true

    provisioner "remote-exec" {
            inline = <<CMD
apt-get update && apt-get upgrade -y git

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
