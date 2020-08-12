# cloudflare-backup
cloudflare-backup is a tool that allows automated exporting of all DNS records in a CloudFlare account. This is useful if you want to keep a backup of your domains, or if you need to search through all the records to look for something in particular.

## Usage
You must create a CloudFlare API token first. Follow [these instructions](https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys#12345680), and give the token these permissions at minimum: Zone / DNS / Read and Zone / Zone / Read.

Then, build this program (`go build`) and run it: `./cloudflare-backup -api-token "(your token goes here)"`. DNS records for all of the domains in your account will be exported to `output/`. (you can change this with the `-output` flag)