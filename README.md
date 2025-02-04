# LEAccount

The program assume that an environmental variable LEDir exists.
LEDir point to the directory that holds the LE account files.
After a successful LE account is creates there should exist a LE private and public account key file.  

Naming convention of the key files: <LE account name><Prod|Test><Priv|Pub>.key  
example: azulTestPriv.key and azulTestPub.key

The programs are looking for a file with the name <accountname>_info.yaml. That file is a yaml file that contains contact mailing addresses
for the Let's Encrypt account.



## createLEAccount

Program that that creates an Let's Encrypt Account. the program assumes that a file with the name <account name__info.yaml exists.  
The info file contains the contact information (e.g. mailing addresses).  

usage: ./createLEAccount /acnt=name /type=[prod|test] [/dbg]  

## checkLEAccount

Program that that checks an Let's Encrypt Account.  

usage: ./checkLEAccount /acnt=name /type=[prod|test] [/dbg]  

## LEAccount

A program that combines check and create programs into one.  

usage: ./checkLEAccount /acnt=name [/cmd=[create|check]] /type=[prod|test] [/dbg]  
default is check!  
