import { FreshContext } from "$fresh/server.ts";
import { authCookieName } from "../../lib/auth.ts";

export const handler = (_req: Request, _ctx: FreshContext): Response => {
  console.log("signout");

  const resp = new Response(null, { status: 303 });
  resp.headers.set("location", "/");
  resp.headers.set(
    "Set-Cookie",
    `${authCookieName}=""; expires=${new Date(0).toUTCString()};`,
  );
  return resp;
};
