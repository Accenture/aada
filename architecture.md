Problem Statement
---

Credentials within federated (assumed-role) accounts that don't use AWS SSO are hard to come by.  Accenture grants
access to AWS accounts using a registered Azure application to bounce into an assumed-role via SAML credentials,
however this doesn't easily grant CLI/SDK credentials.  AADA solves this problem in an elegant way.

How does AADA Work?
---

The general architecture is a pair of lambda functions fronted by a pair of API gateways.  This is replicated
into two regions.

```
                            ┌─────────────┐                           
                            │    AADA     │                           
                            │   Client    │                           
                            └─────────────┘                           
                                                                      
┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐
             us-west-1                          us-east-1             
│                                 ││                                 │
  ┌─────────────┐ ┌─────────────┐    ┌─────────────┐ ┌─────────────┐  
│ │  Websocket  │ │    HTTP     │ ││ │  Websocket  │ │    HTTP     │ │
  │ API Gateway │ │ API Gateway │    │ API Gateway │ │ API Gateway │  
│ ├─────────────┤ ├─────────────┤ ││ ├─────────────┤ ├─────────────┤ │
  │  Websocket  │ │    HTTP     │    │  Websocket  │ │    HTTP     │  
│ │   Lambda    │ │   Lambda    │ ││ │   Lambda    │ │   Lambda    │ │
  └─────────────┘ └─────────────┘    └─────────────┘ └─────────────┘  
│         ┌─────────────┐         ││         ┌─────────────┐         │
          │   AWS KMS   │                    │   AWS KMS   │          
│         │   Replica   │         ││         │     Key     │         │
          └─────────────┘                    └─────────────┘          
└ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘└ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘
```

The AADA client is installed as a 
[credential process](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html)
which effectively calls AADA whenever credentials are required for an API call.  AADA uses three main components
to coordinate getting credentials distributed.  From a high level, the client initiates a websocket connection
to get a session started.  It then uses that session to launch an Azure OIDC authentication flow, passing the
session in the state parameter of the 
[OIDC auth flow](https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-protocols-oidc)
so it will be posted to the HTTP lambda callback.  The HTTP lambda validates the auth token, then validates the
session state, then checks to see if the user belongs to the group they're asking for credentials to.  Assuming
all of that passes, the HTTP lambda does an 
[STS AssumeRole](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)
using the AADA trustpoint into the target role.  With assumed role credentials, the HTTP lambda then sends those
back to the client via the initial Websocket connection.

Visually, the authentication flow looks like this:

```
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│    AADA     │ │  Websocket  │ │    HTTP     │ │  Azure AD   │
│   Client    │ │   Lambda    │ │   Lambda    │ │    OIDC     │
└──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └──────┬──────┘
       │               │               │               │       
       ├──────────────▶│               │               │       
       │  New Session  │               │               │       
       │               │               │               │       
       │◀──────────────┤               │               │       
       │ Signed Session│               │               │       
       │               │               │               │       
       ├───────────────┼───────────────┼──────────────▶│       
       │              Web Authentication               │       
       │               │               │               │       
       │               │               │◀──────────────┤       
       │               │               │ Token Callback│       
       │               │               │               │       
       │◀──────────────┤◀──────────────┤               │       
       │  Credentials  │  Credentials  │               │       
       │               │               │               │       
┌──────┴──────┐ ┌──────┴──────┐ ┌──────┴──────┐ ┌──────┴──────┐
│    AADA     │ │  Websocket  │ │    HTTP     │ │  Azure AD   │
│   Client    │ │   Lambda    │ │   Lambda    │ │    OIDC     │
└─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘
```
