"use client";

import { signIn, signOut } from "next-auth/react";
import Link from "next/link";

export const LoginButton = () => {
  return (
    <button onClick={() => signIn("asgardeo")}>
      Sign in
    </button>
  );
};


export const LogOutButton = () => {
  return (
    <button onClick={() => signOut()}>
      Sign out
    </button>
  );
};

