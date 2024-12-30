# LEAccount

The program assume that an environmental variable LEDir exists.
LEDir point to the directory that holds the LE account files.
After a successful LE account is creates there should exist a LE private and public account key file.  

Naming convention: <LE account name><Prod|Test><Priv|Pub>.key  
example: azulTestPriv.key and azulTestPub.key

## createLEAccount

Program that that creates an Let's Encrypt Account. the program assumes that a file with the name <account name__info.yaml exists.  
The info file contains the contact information (e.g. mailing addresses).  

usage: ./createLEAccount /acnt=name /type=[prod|test] [/dbg]  

## check LEAccount

Program that that checks an Let's Encrypt Acoount.  

usage: ./checkLEAccount /acnt=name /type=[prod|test] [/dbg]  

