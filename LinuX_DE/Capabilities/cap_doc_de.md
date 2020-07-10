

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
