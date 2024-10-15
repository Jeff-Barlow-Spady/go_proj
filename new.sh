########################################################
# System installer Fedora Dev machine 
# Date: 10-10-2024
# OS: Fedora
# Author: Rutger Laurman
########################################################

echo "########################################################"
echo ""
echo "Start automated system install"
echo ""
echo "Based on Omakub installer"
echo ""
echo "########################################################"

################
# Update system
################
sudo dnf update -y
sudo dnf upgrade -y
sudo dnf install -y curl git unzip zsh

# Switch to zshell
chsh -s $(which zsh)

touch ~/.zshrc
echo 'export PATH="$HOME/bin:$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

########################################################
# Terminal
########################################################

################
# Install fastfetch
################
sudo dnf install -y fastfetch

################
# Install neovim
################
cd /tmp
wget -O nvim.tar.gz "https://github.com/neovim/neovim/releases/latest/download/nvim-linux64.tar.gz"
tar -xf nvim.tar.gz
sudo install nvim-linux64/bin/nvim /usr/local/bin/nvim
sudo cp -R nvim-linux64/lib /usr/local/
sudo cp -R nvim-linux64/share /usr/local/
rm -rf nvim-linux64 nvim.tar.gz
cd -

################
# Install lazygit
################
cd /tmp
LAZYGIT_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | grep -Po '"tag_name": "v\K[^"]*')
curl -sLo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_${LAZYGIT_VERSION}_Linux_x86_64.tar.gz"
tar -xf lazygit.tar.gz lazygit
sudo install lazygit /usr/local/bin
rm lazygit.tar.gz lazygit
cd -

################
# Install terminal tools
################
sudo dnf install -y fzf ripgrep bat exa zoxide mlocate btop httpd-tools fd-find tldr

################
# Install Docker
################
# Remove old versions
sudo dnf remove -y docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine

# Set up the repository
sudo dnf -y install dnf-plugins-core
sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo

# Install Docker engine and standard plugins
sudo dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Start Docker and enable on boot
sudo systemctl start docker
sudo systemctl enable docker

# Give this user privileged Docker access
sudo usermod -aG docker ${USER}

# Limit log size to avoid running out of disk
echo '{"log-driver":"json-file","log-opts":{"max-size":"10m","max-file":"5"}}' | sudo tee /etc/docker/daemon.json

################
# Install libraries
################
sudo dnf groupinstall -y "Development Tools"
sudo dnf install -y \
  pkgconfig autoconf bison clang rust \
  openssl-devel readline-devel zlib-devel libyaml-devel ncurses-devel libffi-devel gdbm-devel jemalloc-devel \
  vips ImageMagick ImageMagick-devel mupdf mupdf-tools libgtop2-devel clutter-devel \
  redis sqlite sqlite-devel mariadb-devel postgresql-devel postgresql

################
# Install mise
################
sudo dnf install -y gnupg wget curl
curl https://get.mise.sh | sh

################
# Install PHP + Composer
################
sudo dnf install -y https://rpms.remirepo.net/fedora/remi-release-$(rpm -E %{fedora}).rpm
sudo dnf module reset php
sudo dnf module install php:remi-8.3
sudo dnf install -y php php-cli php-common php-{curl,intl,mbstring,opcache,pdo_mysql,pdo_pgsql,sqlite3,redis,xml,zip}

php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
php composer-setup.php --quiet && sudo mv composer.phar /usr/local/bin/composer
rm composer-setup.php

################
# Install Node.js, Go, Rust, and Python using mise
################

# Install necessary plugins
mise plugin install node
mise plugin install go
mise plugin install rust
mise plugin install python

# Install the latest versions
mise install node latest
mise install go latest
mise install rust latest
mise install python latest

# Set the installed versions as global defaults
mise global node latest
mise global go latest
mise global rust latest
mise global python latest


########################################################
# Desktop
########################################################

################
# Install flatpak
################
sudo dnf install -y flatpak
sudo flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo

################
# Install kitty
################
sudo dnf install -y kitty

################
# Install flameshot
################
sudo dnf install -y flameshot

################
# Install gnome-sushi
################
sudo dnf install -y gnome-sushi

# Install gnome tweak tool
sudo dnf install -y gnome-tweaks

# Install localsend - airdrop alternative
cd /tmp
LOCALSEND_VERSION=$(curl -s "https://api.github.com/repos/localsend/localsend/releases/latest" | grep -Po '"tag_name": "v\K[^"]*')
wget -O localsend.rpm "https://github.com/localsend/localsend/releases/latest/download/LocalSend-${LOCALSEND_VERSION}-linux-x86-64.rpm"
sudo dnf install -y ./localsend.rpm
rm localsend.rpm
cd -

# Install Obsidian
flatpak install -y flathub md.obsidian.Obsidian

# Install Pinta - photoshop alternative
flatpak install -y flathub com.github.PintaProject.Pinta

# Install VLC
sudo dnf install -y vlc

# Install VSCode
sudo rpm --import https://packages.microsoft.com/keys/microsoft.asc
sudo sh -c 'echo -e "[code] \nname=Visual Studio Code \nbaseurl=https://packages.microsoft.com/yumrepos/vscode \nenabled=1 \ngpgcheck=1 \ngpgkey=https://packages.microsoft.com/keys/microsoft.asc" > /etc/yum.repos.d/vscode.repo'
sudo dnf check-update
sudo dnf install -y code

# Install xournalpp - sign PDF documents
sudo dnf install -y xournalpp

# Install ulauncher
flatpak install -y flathub io.ulauncher.ULAUNCHER

# Install 1Password
sudo rpm --import https://downloads.1password.com/linux/keys/1password.asc
sudo sh -c 'echo -e "[1password] \nname=1Password Channel \nbaseurl=https://downloads.1password.com/linux/rpm/stable/\$basearch \nenabled=1 \ngpgcheck=1 \ngpgkey=https://downloads.1password.com/linux/keys/1password.asc" > /etc/yum.repos.d/1password.repo'
sudo dnf install -y 1password 1password-cli

######################################
# Install Signal
######################################
curl -s https://updates.signal.org/desktop/rpm/signing_key.asc | sudo rpm --import -
sudo sh -c 'echo -e "[signal] \nname=Signal Desktop \nbaseurl=https://updates.signal.org/desktop/rpm/x86_64/ \nenabled=1 \ngpgcheck=1 \ngpgkey=https://updates.signal.org/desktop/rpm/signing_key.asc" > /etc/yum.repos.d/signal.repo'
sudo dnf install -y signal-desktop

############## 
# Set Gnome Extensions
############## 

sudo dnf install -y gnome-extensions-app pipx
pipx install gnome-extensions-cli --system-site-packages
pipx ensurepath

# Install new extensions
gext install tactile@lundal.io
gext install just-perfection-desktop@just-perfection
gext install blur-my-shell@aunetx
gext install space-bar@luchrioh
gext install undecorate@sun.wxg@gmail.com
gext install tophat@fflewddur.github.io
gext install AlphabeticalAppGrid@stuarthayhurst

# Compile gsettings schemas in order to be able to set them
sudo cp ~/.local/share/gnome-shell/extensions/tactile@lundal.io/schemas/org.gnome.shell.extensions.tactile.gschema.xml /usr/share/glib-2.0/schemas/
sudo cp ~/.local/share/gnome-shell/extensions/just-perfection-desktop@just-perfection/schemas/org.gnome.shell.extensions.just-perfection.gschema.xml /usr/share/glib-2.0/schemas/
sudo cp ~/.local/share/gnome-shell/extensions/blur-my-shell@aunetx/schemas/org.gnome.shell.extensions.blur-my-shell.gschema.xml /usr/share/glib-2.0/schemas/
sudo cp ~/.local/share/gnome-shell/extensions/space-bar@luchrioh/schemas/org.gnome.shell.extensions.space-bar.gschema.xml /usr/share/glib-2.0/schemas/
sudo cp ~/.local/share/gnome-shell/extensions/tophat@fflewddur.github.io/schemas/org.gnome.shell.extensions.tophat.gschema.xml /usr/share/glib-2.0/schemas/
sudo cp ~/.local/share/gnome-shell/extensions/AlphabeticalAppGrid@stuarthayhurst/schemas/org.gnome.shell.extensions.AlphabeticalAppGrid.gschema.xml /usr/share/glib-2.0/schemas/
sudo glib-compile-schemas /usr/share/glib-2.0/schemas/

# Configure Tactile
gsettings set org.gnome.shell.extensions.tactile col-0 1
gsettings set org.gnome.shell.extensions.tactile col-1 2
gsettings set org.gnome.shell.extensions.tactile col-2 1
gsettings set org.gnome.shell.extensions.tactile col-3 0
gsettings set org.gnome.shell.extensions.tactile row-0 1
gsettings set org.gnome.shell.extensions.tactile row-1 1
gsettings set org.gnome.shell.extensions.tactile gap-size 32

# Configure Just Perfection
gsettings set org.gnome.shell.extensions.just-perfection animation 2
gsettings set org.gnome.shell.extensions.just-perfection dash-app-running true
gsettings set org.gnome.shell.extensions.just-perfection workspace true
gsettings set org.gnome.shell.extensions.just-perfection workspace-popup false

# Configure Blur My Shell
gsettings set org.gnome.shell.extensions.blur-my-shell.appfolder blur false
gsettings set org.gnome.shell.extensions.blur-my-shell.lockscreen blur false
gsettings set org.gnome.shell.extensions.blur-my-shell.screenshot blur false
gsettings set org.gnome.shell.extensions.blur-my-shell.window-list blur false
gsettings set org.gnome.shell.extensions.blur-my-shell.panel blur false
gsettings set org.gnome.shell.extensions.blur-my-shell.overview blur true
gsettings set org.gnome.shell.extensions.blur-my-shell.overview pipeline 'pipeline_default'
gsettings set org.gnome.shell.extensions.blur-my-shell.dash-to-dock blur true
gsettings set org.gnome.shell.extensions.blur-my-shell.dash-to-dock brightness 0.6
gsettings set org.gnome.shell.extensions.blur-my-shell.dash-to-dock sigma 30
gsettings set org.gnome.shell.extensions.blur-my-shell.dash-to-dock static-blur true
gsettings set org.gnome.shell.extensions.blur-my-shell.dash-to-dock style-dash-to-dock 0

# Configure Space Bar
gsettings set org.gnome.shell.extensions.space-bar.behavior smart-workspace-names false
gsettings set org.gnome.shell.extensions.space-bar.shortcuts enable-activate-workspace-shortcuts false
gsettings set org.gnome.shell.extensions.space-bar.shortcuts enable-move-to-workspace-shortcuts true
gsettings set org.gnome.shell.extensions.space-bar.shortcuts open-menu "@as []"

# Configure TopHat
gsettings set org.gnome.shell.extensions.tophat show-icons false
gsettings set org.gnome.shell.extensions.tophat show-cpu false
gsettings set org.gnome.shell.extensions.tophat show-disk false
gsettings set org.gnome.shell.extensions.tophat show-mem false
gsettings set org.gnome.shell.extensions.tophat network-usage-unit bits

# Configure AlphabeticalAppGrid
gsettings set org.gnome.shell.extensions.alphabetical-app-grid folder-order-position 'end'

############################ 
# Set Gnome keyboard shortcuts
############################ 

# Alt+F4 is very cumbersome
gsettings set org.gnome.desktop.wm.keybindings close "['<Super>w']"

# Make it easy to maximize like you can fill left/right
gsettings set org.gnome.desktop.wm.keybindings maximize "['<Super>Up']"

# Make it easy to resize undecorated windows
gsettings set org.gnome.desktop.wm.keybindings begin-resize "['<Super>BackSpace']"

# Use super for workspaces
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-1 "['<Super>1']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-2 "['<Super>2']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-3 "['<Super>3']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-4 "['<Super>4']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-5 "['<Super>5']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-6 "['<Super>6']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-7 "['<Super>7']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-8 "['<Super>8']"
gsettings set org.gnome.desktop.wm.keybindings switch-to-workspace-9 "['<Super>9']"

# Reserve slots for custom keybindings
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings \
 "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/', \
   '/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom1/', \
   '/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom2/', \
   '/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom3/']"

# Set ulauncher to Super+Space
gsettings set org.gnome.desktop.wm.keybindings switch-input-source "@as []"
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ name 'ulauncher-toggle'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ command 'ulauncher-toggle'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ binding '<Super>space'

# Set flameshot on alternate print screen key
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom1/ name 'Flameshot'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom1/ command 'sh -c -- "flameshot gui"'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom1/ binding '<Control>Print'

############################ 
# Set Gnome settings
############################ 
# Center new windows in the middle of the screen
gsettings set org.gnome.mutter center-new-windows true

# Set monospace font (adjust as needed)
# gsettings set org.gnome.desktop.interface monospace-font-name 'CaskaydiaMono Nerd Font 10'

# Reveal week numbers in the Gnome calendar
gsettings set org.gnome.desktop.calendar show-weekdate true

############################ 
# Set Gnome theme
############################ 

gsettings set org.gnome.desktop.interface color-scheme 'prefer-dark'

############################ 
# Change shell to zsh and install and load starship
############################ 

curl -sS https://starship.rs/install.sh | sh
echo 'eval "$(starship init zsh)"' >> ~/.zshrc
source ~/.zshrc

echo "########################################################"
echo "Installation complete. Please reboot!"
echo "########################################################"
