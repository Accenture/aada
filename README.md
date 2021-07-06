# AADA
Accenture Active Directory Authenticator

## Release information
Release 1.0.1-beta, published July 6, 2021 - A few minor tweaks and the
addition of full autoconfiguration, assuming `aada` is in your path.

Release 1.0.0-beta, published June 9, 2021 - Complete overhaul and 
refactor to use OAuth2 as a trusted Azure AD application.  This version
is 100% different, and not backward compatible with prior versions.  It 
requires a custom trust on the assumed role, but with that trust, buys
a completely seamless CLI authentication experience.

Release 0.1.8, published March 15, 2021 - New builds with Go version
1.16.2.

Release 0.1.7, published Feb 26, 2021 - New builds with Go version 1.16.0
to fix the false-positive hueristic virus detection by Cylance and
Microsoft Defender (Program.Win32/Wacapew.C!ml).  To be clear, the 
previous version had no virus in it.  This version is actually throwing
a false positive for TrojanSpy.MSIL.bgkz with the Jiangmin scanner, but
Defender is now happy.  I'll try to run these through VirusTotal from
now on to make sure we don't have any more false positives.

Release 0.1.6, published Feb 22, 2021 - Added the ability for AADA to
automatically create the .aws folder if it doesn't already exist.

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
This is an AWS SDK credential helper that understands how to authenticate
you against the Accenture Active Directory, providing a completely seamless
authentication experience without having to enter your password or security
token.  This is useful for tools such as Terraform to more 
easily work within Accenture AWS accounts using federated credentials.

## How do I install it?
There is a binaries folder that includes binaries for Linux, Mac, and 
Windows.  Download the appropriate binary and place it into your path.
There are no other system requirements.  The "aada" binary must be in
your system path for this to work properly.

## How do I use it?
Start with configuration by running `aada -configure` and aada will setup
your granted profiles into your AWS configuration.  Each profile will be
configured to automatically call aada with the correct account and group
information by any app that uses the AWS SDK.

To test it out, run trusty get-caller-identity.  On the first run, the 
aada authentication pop-through should come up, and your CLI call should
complete successfully.
```
$ aws --profile iesawsna-sandbox sts get-caller-identity
{
    "UserId": "AROA4UGSQ27FZ4TPYZFLZ:eric.hill@accenture.com",
    "Account": "868024899531",
    "Arn": "arn:aws:sts::868024899531:assumed-role/iesawsna-sandbox/eric.hill@accenture.com"
}
```

Once your first authentication completes, the credentials are cached in the
aws credentials file so that subsequent API calls complete without the 
authentication pop-through.  Credentials are good for an hour by default,
and with a completely transparent experience, support for longer assumption
times is not currently planned.  Please reach out with good use-cases for
longer assumption time.

## What did that just do?
1. Check for cached credentials and just return them if they're still valid
2. Initiate a cli login session using a websocket to wss.aabg.io
3. Initiate an OAuth2 authentication against the aada application, with a redirect to aabg.io/authenticator
4. The lambda on the back-end verifies the OAuth2 token and validates group membership
5. The lambda then retrieves assumed-role credentials and sends them back down the websocket to the client
6. The client caches the credentials and returns them to the SDK

## So what's required on the assumed roles?

Any role that aada needs to assume must trust `arn:aws:iam::464079168809:role/aada-trustpoint` to assume it.
Without this trust, aada cannot give you credentials.  For common shared accounts (the sandbox), this is 
already done.  For other accounts you might be using, the role may have to be updated.

Further, there is a very specific group format that aada uses.  AWS-\[AWS account number\]-\[Role name here\].  The
groups match up to the structure that ACP uses internally.  When you request credentials to one of these roles,
your membership in the Azure AD group is verified before credentials are granted.

## How is aada integrated into the aws config?

Under `~/.aws/config` you'll find profile entries for each role you have access to.

```
[profile iesawsna-sandbox]
credential_process = aada AWS-868024899531-iesawsna-sandbox
```

When you run a command that references this profile (e.g. `aws --profile iesawsna-sandbox sts get-caller-identity`),
the AWS SDK sees the `credential_process` setting and launches `aada AWS-868024899531-iesawsna-sandbox` to get
credentials.

Assuming credentials are received, they are placed into `~/.aws/credentials` for cache purposes.

```
[AWS-868024899531-iesawsna-sandbox]
aws_access_key_id     = ...
aws_secret_access_key = ...
aws_session_token     = ...
expiration_date       = 2021-06-10T15:45:28Z
```

## Who do I blame when things go wrong?
This was written by Eric Hill.  Ping me and I'll see what I can do to help.
