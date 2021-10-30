# GOnk - a Discord bot written in Go
GOnk is named for [Gonk Droids](https://starwars.fandom.com/wiki/GNK_power_droid).

### About
GOnk is a Discord bot whose primary purpose is to listen to the event stream API and track when users in your Discord start streaming on Twitch.
This functionality works if they have Twitch integrated with their Discord account. It will post to a Discord channel of your choosing with a message
that the user has gone live.

### Maintainers
Brandon Barrow

### Usage

### Installation

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add gonk-repo https://brandonlbarrow.github.io/gonk

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
gonk` to see the charts.

To install the gonk chart:

    helm install gonk bbarrow/gonk

To uninstall the chart:

    helm delete gonk