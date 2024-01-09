---
layout: ""
page_title: "icdc provider configuration"
description: |-
    
---

# ICDC Provider
## Schema
### Required
- `auth_group` (String, Sensitive) - User active group, contains needed account and role
- `username` (String, Sensitive)
- `location` (String, Sensitive) - operated location

### Optional
- `password` (String, Sensitive) - user password, also user can declare it using env variable `ICDC_PASSWORD`
- `sso_client_id` (String, Sensitive)
- `sso_realm` (String, Sensitive) - basically operator name 
- `sso_url` (String, Sensitive)



