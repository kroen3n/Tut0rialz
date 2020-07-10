

<b> Linux Capabilities und D0cker-Container </b>


       Für   den   Zweck   der   Durchführung  von  Rechteprüfungen  unterscheiden  traditionelle
       UNIX-Implementierungen zwei Arten von Prozessen: Privilegierte Prozesse  (deren  effektive
       Benutzer-ID  0  ist,  auch  als  Superuser oder Root benannt) und unprivilegierte Prozesse
       (deren effektive UID von Null verschieden  ist).  Privilegierte  Prozesse  übergehen  alle
       Kernel-Rechteprüfungen,   während   unprivilegierte  Prozesse  der  vollen  Rechteprüfung,
       basierend auf den Berechtigungsnachweisen des  Prozesses  (normalerweise:  effektive  UID,
       effektive GID und ergänzende Gruppenliste), unterliegen.

       Beginnend  mit  Kernel  2.2  unterteilt  Linux  die  Privilegien, die traditionell mit dem
       Superuser assoziiert sind, in getrennte Einheiten,  die  als  Capabilities  bekannt  sind.
       Diese  können  unabhängig voneinander aktiviert oder deaktiviert werden. Capabilities sind
       ein Attribut pro Thread.

Zitat aus: http://manpages.ubuntu.com/manpages/bionic/de/man7/capabilities.7.html

<i> Man soll es lesen, bevor mit diesem Tutorial beginnt. </i> 
<br> 
</br>
<i> ...und die folgenden Links: </i> 

<i> https://man7.org/linux/man-pages/man8/setcap.8.html </i>

<i> https://man7.org/linux/man-pages/man8/getcap.8.html </i>

<i> https://man7.org/linux/man-pages/man1/capsh.1.html </i> 


<br></br>

<b> Das erste Beispiel - Linux-Host   </b> 

Der Linux-Kernel implementiert eine Vielzahl von Fähigkeiten; <br>
Das folgende Beispiel gibt einen kleinen Überblick darüber, warum und wie man die Macht der Fähigkeiten nutzt. </br>


 Benutzer "hue" und Gruppe "hue" erstellen 
```
root@kroen3n:/home/hue# useradd hue
root@kroen3n:/home/hue# mkdir -p /home/hue
root@kroen3n:/home/hue# chown -R hue:hue /home/hue
```

Prüfen durchführen

```
root@kroen3n:/home/hue# cat /etc/passwd | grep hue
hue:x:1000:1000::/home/hue:/bin/sh
root@kroen3n:/home/hue# 
root@kroen3n:/home/hue# su - hue
$ bash
hue@kroen3n:~$ pwd
/home/hue
```

Eine leere Datei als Root-Benutzer erstellen.
Speicherort bleibt unter dem Homeverzeichnis des hue-Benutzers.

```
root@kroen3n:/home/hue# touch hiya.txt
root@kroen3n:/home/hue#
root@kroen3n:/home/hue# ls -ltr hiya*
-rw-r--r-- 1 root root    0 Jul  9 13:18 hiya.txt
```

 "hue"-Benutzer zu werden, und Operationen auf die Datei anwenden 
 
 ```
 root@kroen3n:/home/hue# su - hue
$ bash
hue@kroen3n:~$ ls -ltr hiya*
-rw-r--r-- 1 root root    0 Jul  9 13:18 hiya.txt
```
Man wird versuchen, in diese Datei zu schreiben; Ich werde ein Golang-Programm ausführen, <a href="https://raw.githubusercontent.com/kroen3n/Tut0rialz/master/LinuX/Capabilities/write_into_file.go"> write_into_file.go</a>, um ein paar Zeilen hinzuzufügen und anzuhängen.

```
package main

import (
  "os"
	"io/ioutil"
	"log"
)

func main(){
	err := ioutil.WriteFile(os.Args[1], []byte("hiya\n"), 0644)

	if err != nil{
		log.Fatal(err)
	}

	file, err := os.OpenFile(os.Args[1], os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil{
		log.Println(err)
	}

	defer file.Close()

	if _, err := file.WriteString("hiya again\n"); err != nil{
		log.Fatal(err)
	}

}
```

Das Programm mit einem Argument ausführen:

```
hue@kroen3n:~$ go run write_into_file.go hiya.txt 
2020/07/09 13:24:15 open hiya.txt: permission denied
exit status 1
```
Hmm... Keine Berechtigung...

Es ist an der Zeit, dass ich etwas anderes versuche.   Ich werde versuchen, die Datei umzubenennen ...
Man könnte den Befehl "mv" verwenden, aber sollte man das jetzt tun? Lassen Sie uns mehr Golang üben, mit dem Programm rename_me.go:

```
package main

import (
	"log"
	"os"
)

func main(){
	actualFile:= os.Args[1]
	newFile := os.Args[2]

	err := os.Rename(actualFile, newFile)

	if err != nil {
		log.Fatal(err)
	}
}
```





```
package main

import (
	"log"
	"os"
)

func main(){
	actualFile:= os.Args[1]
	newFile := os.Args[2]

	err := os.Rename(actualFile, newFile)

	if err != nil {
		log.Fatal(err)
	}
}
```
  
   