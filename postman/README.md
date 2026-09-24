# Nehmya API requests

Import `Nehmya API.postman_collection.json` and `environments/Nehmya Local.postman_environment.json` into Postman, then select **Nehmya Local**. Set `baseUrl` to the Go server address if it differs from `http://localhost:8080`.

## Typical flow

1. Run **Auth → Register user** once, or use an existing account.
2. Run **Auth → Login**. Postman retains the session cookie for subsequent requests to the same host.
3. Run **Tenants → Register tenant**, then **Tenants → List my tenants**. Registration provisions a separate, initially empty PostgreSQL database.
4. For subscription date changes, log in as a user with the `sys_admin` role, set `tenantId` from the tenant response, then run **Subscription Admin → Set subscription dates**.

The registration request creates a tenant and database each time it runs. The admin subscription request expects RFC3339 timestamps. The default subscription is 30 days and can be changed with `DEFAULT_SUBSCRIPTION_DAYS`. Database provisioning uses `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and optional `DB_ADMIN_DATABASE` (defaults to `postgres`); the PostgreSQL role must have `CREATEDB` permission.
