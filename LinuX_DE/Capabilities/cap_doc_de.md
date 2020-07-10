

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

Durchführung der Prüfung

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
Man wird versuchen, in diese Datei zu schreiben; Ich werde ein Golang-Programm ausführen,  um ein paar Zeilen hinzuzufügen und anzuhängen.

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

Man soll zwei Argumente für dieses Programm angeben - den Namen der eigentlichen Datei und den Namen, den man wählt:

```
hue@kroen3n:~$  go run rename_me.go hiya.txt hielau.txt
hue@kroen3n:~$  
hue@kroen3n:~$  ls -ltr hielau.*
-rw-r--r-- 1 root root 0 Jul  9 13:31 hielau.txt
hue@kroen3n:~$  
hue@kroen3n:~$ 
```

Das hat funktioniert...

Und zurück zu meiner (jetzt umbenannten) Datei:

```
hue@kroen3n:~$  ls -ltr hielau.*
-rw-r--r-- 1 root root 0 Jul  9 13:31 hielau.txt
```
Wie man gut erkennen kann, versuche ich, in eine Datei zu schreiben, die dem Benutzer root gehört, während ich ein Nicht-Root-Benutzer bin:

```
hue@kroen3n:~$  id
uid=1000(hue) gid=1000(hue) groups=1000(hue)
```

Dies scheinen Eigentumsfragen zu sein.
Welche Optionen habe ich? - lassen Sie uns über "Fähigkeiten" nachdenken. In der Dokumentation wird erwähnt:

```
      CAP_CHOWN
              beliebige Änderungen an Datei-UIDs und GIDs vornehmen (siehe chown(2))
```	    

"chown" suchen:

```
hue@kroen3n:~$ whereis chown
chown: /bin/chown
hue@kroen3n:~$
```

Kopieren Sie /bin/chown in das Heimatverzeichnis des Benutzers "hue", /home/hue:

```
hue@kroen3n:~$ cp /bin/chown .
hue@kroen3n:~$ ls -ltr /home/hue/chown
-rwxr-xr-x 1 hue  hue  72512 Jul  9 14:20 chown
```

Schauen wir nach, wo sich die Tools getcap und setcap befinden. 

```
hue@kroen3n:~$ whereis getcap
getcap: /sbin/getcap
hue@kroen3n:~$ 
hue@kroen3n:~$ whereis setcap
setcap: /sbin/setcap
hue@kroen3n:~$ 
```

```
hue@kroen3n:~$  
hue@kroen3n:~$  /sbin/getcap /home/hue/chown
hue@kroen3n:~$ 
```

Wie erwartet, keine Ausgabe ...

CAP_CHOWN mit setcap aktivieren:

```
hue@kroen3n:~$  /sbin/setcap cap_chown+ep chown 
unable to set CAP_SETFCAP effective capability: Operation not permitted
hue@kroen3n:~$ 
```

<i> Man sollte sudo verwenden oder root werden, um diesen Fehler zu überspringen. </i>



```
hue@kroen3n:~$  sudo /sbin/setcap cap_chown+ep chown 
hue@kroen3n:~$  
hue@kroen3n:~$ pwd
/home/hue
hue@kroen3n:~$
hue@kroen3n:~$ /sbin/getcap ./chown
chown = cap_chown+ep
```
Ausführung von Kommando:

```
hue@kroen3n:~$ 
hue@kroen3n:~$ ./chown hue:hue hielau.txt
hue@kroen3n:~$
hue@kroen3n:~$ ls -ltr hiya.py 
-rw-r--r-- 1 hue hue 0 Jul  9 14:41 hielau.txt
```

Jetzt werde ich versuchen, etwas in meine Datei zu schreiben: 

```
hue@kroen3n:~$ go run write_into_file.go hielau.txt 
hue@kroen3n:~$ more hielau.txt 
hiya
hiya again
hue@kroen3n:~$
```




