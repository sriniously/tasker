import "./index.css";
import { AppLayout } from "@/components/layouts/app-layout";
import { AuthLayout } from "@/components/layouts/auth-layout";
import { ProtectedRoute } from "@/components/protected-route";
import { PublicRoute } from "@/components/public-route";
import { CategoriesPage } from "@/pages/categories-page";
import { DashboardPage } from "@/pages/dashboard-page";
import { LandingPage } from "@/pages/landing-page";
import { SettingsPage } from "@/pages/settings-page";
import { TodosPage } from "@/pages/todos-page";
import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
} from "react-router-dom";

const routes = createRoutesFromElements(
  <>
    <Route
      path="/"
      element={
        <PublicRoute>
          <LandingPage />
        </PublicRoute>
      }
    />
    <Route
      path="/auth/*"
      element={
        <PublicRoute>
          <AuthLayout />
        </PublicRoute>
      }
    />
    <Route
      path="/dashboard"
      element={
        <ProtectedRoute>
          <AppLayout>
            <DashboardPage />
          </AppLayout>
        </ProtectedRoute>
      }
    />
    <Route
      path="/todos"
      element={
        <ProtectedRoute>
          <AppLayout>
            <TodosPage />
          </AppLayout>
        </ProtectedRoute>
      }
    />
    <Route
      path="/categories"
      element={
        <ProtectedRoute>
          <AppLayout>
            <CategoriesPage />
          </AppLayout>
        </ProtectedRoute>
      }
    />
    <Route
      path="/settings"
      element={
        <ProtectedRoute>
          <AppLayout>
            <SettingsPage />
          </AppLayout>
        </ProtectedRoute>
      }
    />
  </>,
);

const router = createBrowserRouter(routes);

export default router;
