# shortcut-runner
## Usage
Program has init script to run as daemon. Init script is wroten for Debian.
Copy init script from misc to your `/etc/init.d/` directory.
Execute `update-rc.d shortcut-runner defaults`

## Configuration
Program can get options from command line and config files (it's located in /etc/shortcut-runner/shortcut-runner.yml by default). Copy of default config you can get in misc dir of this repo.
Key name you can see in `keymaps.go` file.
