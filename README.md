CSAB stands for Crouton Start and Backup. It is a command line tool to manage, start, and auto backup chroots made with crouton.

Installation:
Should be as simple as downloading the latest csab from https://drive.google.com/file/d/0B-k2N_tmdd_wYm9kemlhSlpxQWM/edit?usp=sharing 
Then starting up your shell making sure your in ~/Downloads or wherever you downloaded it and type:

// This makes it executable

sudo chmod +x csab

// Moving it to the above directory

sudo mv ./csab /usr/local/bin

Since I’ve only installed this on my own pixel I’m hoping for feedback.

Usage:
The first time you run csab it will ask if you wish to add a chroot say yes and follow the instructions. After you finish adding your chroot it will ask if you would like to add another if you only have one or don’t want to add them all at first just type no. You can always add more by using the argument config or -c when starting csab like so: csab -c or csab config. This process will create a config file called chrootconfig.toml in your Downloads folder.

After you have configured your chroots it will detect that you have not made a backup yet and will kick off making a backup. After your first time it will only prompt if you would like to make a backup if your last backup is older than 48 hours. It will then start up your linux chroot.

If people want:
I will add more features if I get comments and people want to use it. Ideas I’ve got for improvements include:
Being able to set how long between backups.

Building it for yourself:
Feel free its written in goLang you will need burntsushi’s toml parser https://github.com/BurntSushi/toml oh actually I didn’t know this at first but you can install goLang on chrome Os https://code.google.com/p/go-wiki/wiki/ChromeOS
