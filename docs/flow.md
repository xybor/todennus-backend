# OAuth2 Flow

## Authorization Code Flow (with/without PKCE)

### Step 1. User clicks the login button on ClientApp.

### Step 2. ClientApp redirects user to `/oauth2/authorize`.

### Step 3. `/oauth2/authorize` checks the session to get all authenticated IdPs.

### Step 4. `/oauth2/authorize` gets all acceptable IdPs of this Client, then compare with the authenticated IdPs.

- If there is at least a same IdP, `/oauth2/authorize` redirects user back to ClientApp with Authorization code.
- If there is no same IdP, `/oauth2/authorize` redirects user to `/oauth2/idps` to show the user a list of available IdPs.
  + User chooses one IdP and `/oauth2/idps` redirects user to that IdP's login page.
  + After user logins, IdP makes a request to `/oauth2/login` to notify about the authentication result.
    + If fails, `/oauth2/login` responds `redirect_uri=client-app-callback-with-error`.
    + If success, `/oauth2/login` adds a new authenticated IdP to the current user session. Then responds `redirect_uri=client-app-callback-with-code`.
  + IdP uses the `redirect_uri` in the response to redirect user.

### Step 6. ClientApp receives the response of OAuth2 Provider and handle the result.
