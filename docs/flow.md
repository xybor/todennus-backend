# OAuth2 Flow

## Current Flow

### Step 1. User clicks the login button on ClientApp.

### Step 2. ClientApp redirects user to `/oauth2/authorize`.

### Step 3. `/oauth2/authorize` checks the session to know if the user is authenticated.

- If yes, `/oauth2/authorize` redirects user back to ClientApp with:
  - Code: if the flow is Authorization Code. Save the code into database.
  - Token: if the flow is Implicit.
  - Then go to step 7.
- If no, `/oauth2/authorize` all information in this flow will be stored in a key called authorization id. Then redirects user to the IdP login page with the authorization id.

### Step 4. After user logins, IdP makes a request to `/auth/callback` to notify about the authentication result.
  + `POST /auth/callback` with the following body:
    ```json
    {
      "authorization_id": "<authorization_id>",
      "success": true,
      "user_id": "<user_id>",
      "username": "<username>",
      "display_name": "<display_name>"
    }
    {
      "authorization_id": "<authorization_id>",
      "success": false,
      "error": "invalid_auth",
      "error_description": "invalid username or password"
    }
    ```
  + `/auth/callback` stores the above information into storage by a key called auth id then responds this id.
  + IdP redirects the user to the `/session/update?auth_id=xxx`.

### Step 5. `/session/update` checks the result in storage.
  + If success, update the session of user.
  + Redirect user to `/oauth2/authorize` with queries are retrieved from information in authorization id.

### Step 6. Go to step 3 with sucess case.

### Step 7. ClientApp receives the response and handle the result.



## Final Flow

### Step 1. User clicks the login button on ClientApp.

### Step 2. ClientApp redirects user to `/oauth2/authorize`.

### Step 3. `/oauth2/authorize` checks the session to get all authenticated IdPs.

### Step 4. `/oauth2/authorize` gets all acceptable IdPs of this Client, then compare with the authenticated IdPs.

- If there is at least a same IdP, `/oauth2/authorize` redirects user back to ClientApp with Authorization code.
- If there is no same IdP, `/oauth2/authorize` redirects user to `/oauth2/idps` to show the user a list of available IdPs.
- User chooses one IdP and `/oauth2/idps` redirects user to that IdP's login page.

### Step 5. After user logins, IdP makes a request to `/oauth2/login` to notify about the authentication result.
  + If fails, `/oauth2/login` responds `redirect_uri=client-app-callback-with-error`.
  + If success, `/oauth2/login` adds a new authenticated IdP to the current user session. Then responds `redirect_uri=client-app-callback-with-code`.

### Step 6. IdP uses the `redirect_uri` in the response to redirect user.

### Step 7. ClientApp receives the response of OAuth2 Provider and handle the result.
