# AADA
Accenture Active Directory Authenticator

## Release information
Release 0.1.5, published Jan 19, 2021 - Added the ability for AADA to 
automatically create config and credentials files if they do not already
exist.  Additionally, the aada Mac binary has been signed with an Apple
developer certificate to comply with the new Mac signature requirements.
The signed binary is in binaries/mac/aada, and has also been zipped into
binaries/mac/aada.zip for your convenience.

Release 0.1.4, published Jan 12, 2021 - No new features.  New Go compiler
producing a slightly better binary and refreshed dependencies.  Not a 
required upgrade.

Initial release 0.1.3, published May 12, 2020.

## What is this?
This tool uses your Accenture enterprise id (your.email@accenture.com) to 
request a SAML token from Active Directory that is exchanged for a set of
AWS credentials to enable both the AWS CLI, as well as applications written
against the AWS SDK.  This is useful for tools such as Terraform to more 
easily work within Accenture AWS accounts using federated credentials.

## How do I install it?
There is a binaries folder that includes binaries for Linux, Mac, and 
Windows.  Download the appropriate binary and place it into your path.
There are no other system requirements.

## How do I use it?
Run aada, enter your enterprise id, password, and current Symantec VIP
token.

```
$ aada
Username: eric.hill@accenture.com
Password: <redacted>
Symantec VIP: 123456
0/5 preparing
1/5 getting new session
2/5 authenticating
3/5 fetching SAML token
You may assume one of the following roles:
   0: arn:aws:iam::012345678901:role/sample-1
   1: arn:aws:iam::868024899531:role/iesawsna-sandbox
   2: arn:aws:iam::012345678901:role/sample-2
Role Number: 1
4/5 exchanging SAML token for role credentials
5/5 installing access key ASIA4UGSQ27FTWSETYAV into profile iesawsna-sandbox
complete - your access expires in 59m59s
```
To confirm the credentials are valid, make a call to get-caller-identity.
```
$ aws --profile iesawsna-sandbox sts get-caller-identity
{
    "UserId": "AROA4UGSQ27FZ4TPYZFLZ:eric.hill",
    "Account": "868024899531",
    "Arn": "arn:aws:sts::868024899531:assumed-role/iesawsna-sandbox/eric.hill"
}
```

## What did that just do?
1. Exchanged your information for a SAML token
2. Exchanged the SAML token for AWS credentials
3. Placed those credentials in ~/.aws/credentials under a profile named after the role
4. Ensured ~/.aws/config contains both an output type and region for that profile

## Options and switches
```
$ aada -help
Usage of ./binaries/mac/aada:
  -d duration
        duration of assumed credentials (default 1h0m0s)
  -trace string
        if specified, traces all activity into this file
  -u string
        specify the username
  -version
        display version information and exit
```
The most useful switch is the -d duration switch.  For the AABG sandbox, you'll
probably want:
```
$ aada -d 8h
```

## Who do I blame when things go wrong?
This was written by Eric Hill.  Don't blame me, just send me a trace and 
description of what broke and I'll see if I can get it fixed.
```
$ aada -trace=tracefile.log
```
Your password will NOT be included in the trace file.

## Future work
I'm waiting on some integration information so I can extend this to use
OpenID instead of form posts, which should allow a seamless authentication
experience without entering pesky passwords or VIP tokens every time.  More
to come.
