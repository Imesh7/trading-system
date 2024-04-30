import NextAuth, { NextAuthOptions } from "next-auth"

export const authOptions: NextAuthOptions = {
    providers: [
        {
            id: "asgardeo",
            name: "Asgardeo",
            clientId: process.env.ASGARDEO_CLIENT_ID,
            clientSecret: process.env.ASGARDEO_CLIENT_SECRET,
            issuer: process.env.ASGARDEO_ISSUER,
            userinfo: `https://api.asgardeo.io/t/${process.env.ASGARDEO_ORGANIZATION_NAME}/oauth2/userinfo`,
            type: "oauth",
            wellKnown: `https://api.asgardeo.io/t/${process.env.ASGARDEO_ORGANIZATION_NAME}/oauth2/token/.well-known/openid-configuration`,
            authorization: {
                params:
                    { scope: "openid profile" }
            },
            idToken: true,
            checks: ["pkce", "state"],
            profile(profile) {                
                return {
                    id: profile.sub,
                    name: profile.name,
                    email: profile.email,
                }
            },
        },
    ],
    secret: process.env.NEXTAUTH_SECRET,
    session: {
        strategy: "jwt",
    },
    callbacks: {},
    debug: true
}

const handler = NextAuth(authOptions);
export { handler as GET, handler as POST };