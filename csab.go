package main

import (
	"os"
	"fmt"
	"regexp"
	"time"
	"os/exec"
	confirm "github.com/aintnorest/csab/confirmation"
	tparse "github.com/BurntSushi/toml"
)
type tomlConfig struct {
	Title	string
	Chroots	map[string]chroot
}

type chroot struct {
	Name			string
	Backup_location	string
	Desktop			string
}

func tomlDecode(location string)(*tomlConfig) {
	var config tomlConfig
	if _, err := tparse.DecodeFile(location, &config); err != nil {
		fmt.Println(err, "Not a Valid TOML")
	}
	return &config
}

func CreateConfig() {
	//makes the Config File
 ConfigTomil, err := os.Create("/home/chronos/user/Downloads/chrootconfig.toml")
 	if err != nil {
 		fmt.Println(err)
 	}
 	if err == nil {
 		defer ConfigTomil.Close()
 	}
//add the title
_, err = ConfigTomil.WriteString("title = \"TOML Configuration for Crouton Start and Backup Automation\"\n\n[chroots]")
	if err != nil {
		fmt.Println(err)
	}
}
func AppendConfig() {
	File, err := os.OpenFile("/home/chronos/user/Downloads/chrootconfig.toml", os.O_RDWR|os.O_APPEND, 0660)
	if err == nil {
		defer File.Close()
	} else {
		fmt.Println(err)
		return
	}
	var yn bool = true
	var config tomlConfig
	if _, err := tparse.DecodeFile("/home/chronos/user/Downloads/chrootconfig.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	
	for i := len(config.Chroots) + 1; yn == true; i++ {
	fmt.Println("Would you like to add a Chroot to your config file?  Y/N\n")
	var input string
	fmt.Scanln(&input)
	yn, err := confirm.ConfirmationPrompt(input)
		if err != nil {
			fmt.Println("Confirmation Error")
		}
		if yn == true {
			//Gather Config File Info into three variables
			fmt.Println("\nAll the following inputs are case sensitive.")
			fmt.Println("\nFirst enter the name of the charoot you are trying to add. If you didn't name the chroot with the -n flag, the default naming convention is the distro's version name. For example raring or wheezy. Please Don't name two chroots a like. \nPlease enter the Chroot's name now:\n")
			var chrootname string
			fmt.Scanln(&chrootname)
			fmt.Println("\nNext enter the full pathway location where you would like your backup stored. A good default is /home/chronos/user/Downloads/\n")
			var chrootbl string
			
			for {
				fmt.Scanln(&chrootbl)
				expr, err := regexp.Compile("[/]$")
				if err == nil {
					argumatch := expr.MatchString(chrootbl)
					if argumatch == false {
						chrootbl = chrootbl + "/"
					}
				}
				ffile, err := os.Open(chrootbl)
				if(err == nil) { 
					ffile.Close()
					break }
				fmt.Println("That wasn't a valid Pathway.\n")
				fmt.Println("Re-enter a valid Pathway.\n")
			}
			
			fmt.Println("\nFinally enter the name of the desktop you installed with your chroot. For most it will probably be unity although I personally use cinnamon.\n")
			var chrootdesktop string
			fmt.Scanln(&chrootdesktop)
			_, err = File.WriteString("\n\n\t[chroots."+fmt.Sprintf("%d",i)+"]\n\tname = \""+chrootname+"\"\n\tbackup_location = \""+chrootbl+"\"\n\tdesktop = \""+chrootdesktop+"\"")
				if err != nil {
			fmt.Println(err)
				}
		}

		if yn == false {
			return
		}
	}
	
}

func durationSince(fpathway string) (Elapsed float64, err error){
	file, err := os.Open(fpathway)
	if err != nil {
		return 0, err
	}
	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	if info != nil {
		creationT := info.ModTime() // This should define a new variable :=
		fElapsed := time.Since(creationT)
		Elapsed = fElapsed.Hours()
	}
	file.Close()
	return Elapsed, err
}


func worker(finishedChan chan struct{}) {
	close(finishedChan)
}

func waiting(finishedChan chan struct{}) {
	for {
		timerChan := time.After(2 * time.Second)
		select {
			case <-finishedChan: {
				return
			}
			case <- timerChan: {
				fmt.Print(".")
			}
		}
	}
}

func getFPFromConfig(config tomlConfig, index int) string{
	chroot := config.Chroots[fmt.Sprintf("%d", index)]
	bl := chroot.Backup_location
	n := chroot.Name
	return fmt.Sprintf("%s%sBackup.tar.gz", bl, n)
}

func main() {

	if len(os.Args) > 1 {
		Arguments := os.Args[1]
		expression, err := regexp.Compile("^([cC][oO][nN][fF][iI][gG]|[-][cC])$")
		if err == nil {
			argmatch := expression.MatchString(Arguments)
			if argmatch == false {
				fmt.Println("CSAB is a tool to startup a crouton created chroot and backup scheduler.\n\n Usage: \n\n\t CSAB \n\t\t Will launch the program. \n CSAB -c or config \n\t\t Will allow you to add chroots after you already have the config file setup. You can also just delete the config file.")
				os.Exit(1)
			}
			if argmatch == true {
				AppendConfig()
			}
		}
	}

	_, errrr := os.Open("/home/chronos/user/Downloads/chrootconfig.toml")
	//if err == nil {
	//	file.Close()
	//}
	if errrr != nil {
		CreateConfig()
		AppendConfig()
	}

	var config tomlConfig
	if _, err := tparse.DecodeFile("/home/chronos/user/Downloads/chrootconfig.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	//if errrr == nil {
	//	file.Close()
	//}
	fpathway := ""
	var ChrootN int
	if len(config.Chroots) == 1 {
		fpathway = getFPFromConfig(config, 1)
		ChrootN = 1
	}
	if len(config.Chroots) > 1 {
		fmt.Println("\nHere are a list of your chroots:")
		for chrootname, chroot := range config.Chroots {
			fmt.Printf("\nChroot:%s - %s\n", chrootname, chroot.Name)
		}
		fmt.Println("\nEnter the number of the chroot you wish to enter:\n")
		fmt.Scanln(&ChrootN)
		for {
			if ChrootN >= 1 && ChrootN <= len(config.Chroots) {
			break
			} 
			fmt.Println("You entered an invalid chroot number. Please re-enter a valid number.\n")
			fmt.Scanln(&ChrootN)
		}
		fpathway = getFPFromConfig(config, ChrootN)
	}
fpathway = getFPFromConfig(config, ChrootN)
	
	var BackupWindow float64 = 48
	Duration, errr := durationSince(fpathway)
		finishedChan := make(chan struct{})
	if errr != nil {
		fmt.Println("A backup file wasn't detected. One will be created now.\n")
		go waiting(finishedChan)
		scn := config.Chroots[fmt.Sprintf("%d", ChrootN)].Name
		arg0 := "sudo"
		arg1 := "bash"
		arg2 := "-c"
		arg3 := fmt.Sprintf("edit-chroot -f %s -b %s", fpathway, scn)
		cmd := exec.Command(arg0, arg1, arg2, arg3)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		worker(finishedChan)
	}

fpathway = getFPFromConfig(config, ChrootN)

	if Duration >= BackupWindow {
finishedChan := make(chan struct{})	
		fmt.Println("\nIt has been", BackupWindow, "hours since your last backup. Would you like to make a backup now?\n")
		var input string
		fmt.Scanln(&input)
		fmt.Println("\n")
		Yn, err := confirm.ConfirmationPrompt(input)
		if err != nil {
			fmt.Println("Confirmation Error")
		}
		if Yn == true {
			go waiting(finishedChan)
			scn := config.Chroots[fmt.Sprintf("%d", ChrootN)].Name
			arg0 := "sudo"
			arg1 := "bash"
			arg2 := "-c"
			arg3 := fmt.Sprintf("edit-chroot -f %s -b %s", fpathway, scn)
			cmd := exec.Command(arg0, arg1, arg2, arg3)
			err = cmd.Run()
			if err != nil {
			fmt.Println(err)
			worker(finishedChan)
			}
		}
	}
scd := config.Chroots[fmt.Sprintf("%d", ChrootN)].Desktop
arg00 := "sudo"
arg01 := "bash"
arg02 := "-c"
arg03 := "start"+scd
cmd2 := exec.Command(arg00, arg01, arg02, arg03)
err := cmd2.Run()
if err !=nil {
	fmt.Println(err)
}


}