--  This file is part of the eliona project.
--  Copyright © 2024 LEICOM iTEC AG. All Rights Reserved.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

create schema if not exists zevvy;

-- Should be editable by eliona frontend.
create table if not exists zevvy.configuration
(
    id                      bigserial primary key,
    root_url                text    not null,
    auth_url_path           text    not null,
    client_id               text    not null,
    client_secret           text    not null,
    device_code             text,
    verification_uri        text,
    verification_uri_expire timestamp with time zone,
    verification_interval   integer,
    access_token            text,
    access_token_expire     timestamp with time zone,
    refresh_token           text,
    refresh_interval        integer not null default 60,
    request_timeout         integer not null default 120,
    active                  boolean          default false,
    enable                  boolean          default false,
    user_id                 text,
    project_id              text
);

create table if not exists zevvy.asset_attribute
(
    config_id          integer                  not null,
    asset_id           integer                  not null,
    subtype            text                     not null,
    attribute_name     text                     not null,
    device_reference   text                     not null,
    register_reference text                     not null,
    latest_ts          timestamp with time zone not null default '1900-01-01 00:00:00',
    primary key (config_id, asset_id, subtype, attribute_name)
);

-- Makes the new objects available for all other init steps
commit;
