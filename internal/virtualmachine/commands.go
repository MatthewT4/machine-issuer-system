package virtualmachine

const (
	Uptime = "uptime -s"
	CPU    = "top -bn1 | grep \"Cpu(s)\" | sed \"s/.*, *\\([0-9.]*\\)%* id.*/\\1/\" | awk '{print 100 - $1}'"
	RAM    = "free | awk 'NR==2{printf \"%.2f\", $3*100/$2 }'; echo"
	MEM    = "df -m / | awk 'NR==2{print $4}'"

	Reboot = "reboot"

	CreateUser     = "sudo adduser --disabled-password --gecos '' %s"
	CreatePassword = "echo 'user:%s' | sudo chpasswd"
	GiveRoot       = "sudo usermod -aG sudo %s"

	LimitRoot = `
	echo '%s ALL=(ALL) ALL, !/usr/sbin/userdel admin, !/usr/sbin/deluser admin, !/usr/sbin/usermod admin, !/usr/sbin/groupmod admin, !/usr/sbin/groupdel admin' | sudo tee -a /etc/sudoers.d/%s
	echo '%s ALL=(ALL) ALL, !/bin/chmod /etc/sudoers.d/%s, !/bin/chown /etc/sudoers.d/%s' | sudo tee -a /etc/sudoers.d/%s
	echo '%s ALL=(ALL) ALL, !/bin/chmod /etc/sudoers.d/%s' | sudo tee -a /etc/sudoers.d/admin
	sudo chown admin:admin /etc/sudoers.d/%s
	sudo chmod 0600 /etc/sudoers.d/%s
	`
)

//add user - []string{fmt.Sprintf(CreateUser, login), fmt.Sprintf(CreatePassword, password), fmt.Sprintf(GiveRoot, login), strings.ReplaceAll(LimitRoot, "%s", login)}
