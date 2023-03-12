import {Navigate} from "react-router-dom";
import React from "react";

interface ProtectedRouteProps {
  authenticated: boolean;
  children: JSX.Element;
}

function ProtectedRoute({authenticated, children}: ProtectedRouteProps) {
  if (!authenticated) {
    return <Navigate to="/" replace/>
  }

  return children;
}

export default ProtectedRoute;