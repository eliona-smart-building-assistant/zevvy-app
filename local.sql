--  This file is part of the eliona project.
--  Copyright © 2023 LEICOM iTEC AG. All Rights Reserved.
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

-- This idempotent script resets the database to a defined state ready for testing.
-- The only thing that remains after testing then are the incremented auto-increment values and app
-- registration (which you can optionally remove as well by uncommenting the last command).

SET SCHEMA 'public';

UPDATE public.eliona_store
SET metadata = COALESCE(metadata::jsonb, '{}') ||
               '{"apiUrl": "http://localhost:3030/v1",
                 "apiSpecificationPath": "/version/openapi.json"}'::jsonb
WHERE app_name = 'zevvy';

