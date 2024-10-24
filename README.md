# Get Unassigned ADE Kandji Macs
This binary will output any ADE-connected Mac in Kandji that is not assigned to a user.

## Installation & Usage
Download the `availmacs` binary from the latest release and run it from the terminal on a Mac. You may need to `chmod +x` the file to make it executable.

## Requirements
You need two environmental variables set:
1. `KANDJI_API_KEY`: a Kandji API key that at a minimum has permissions for the `GET` method on the `List ADE Devices`
2. `KANDJI_API_SUBDOMAIN`: your Kandji API subdomain - for example if your API URL was `somecompany.api.kandji.io`, you would set this variable as `somecompany`

Shortcut: `echo -e "export KANDJI_API_KEY=\"key_here\"\nexport KANDJI_API_SUBDOMAIN=\"subdomain_here\"" >> ~/.zshrc` (Fill in your own info in `key_here` and `subdomain_here`)