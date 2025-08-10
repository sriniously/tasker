import { useTheme } from "@/components/theme-provider";
import { ThemeToggle } from "@/components/theme-toggle";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useUser } from "@clerk/clerk-react";
import {
  User,
  Palette,
  Bell,
  Shield,
  Database,
  ExternalLink,
} from "lucide-react";

export function SettingsPage() {
  const { user } = useUser();
  const { theme } = useTheme();

  const themeName =
    theme === "system" ? "System" : theme === "dark" ? "Dark" : "Light";

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground">
          Manage your account and application preferences
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Profile Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="h-5 w-5" />
              Profile
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-4">
              <img
                src={user?.imageUrl}
                alt={user?.fullName || "User"}
                className="w-16 h-16 rounded-full"
              />
              <div>
                <h3 className="font-semibold">{user?.fullName}</h3>
                <p className="text-sm text-muted-foreground">
                  {user?.primaryEmailAddress?.emailAddress}
                </p>
                <Badge variant="secondary" className="mt-1">
                  {(user?.publicMetadata?.role as string) || "User"}
                </Badge>
              </div>
            </div>
            <Separator />
            <div className="space-y-2">
              <p className="text-sm font-medium">Account Details</p>
              <div className="space-y-1 text-sm text-muted-foreground">
                <p>
                  Member since:{" "}
                  {user?.createdAt
                    ? new Date(user.createdAt).toLocaleDateString()
                    : "N/A"}
                </p>
                <p>
                  Last updated:{" "}
                  {user?.updatedAt
                    ? new Date(user.updatedAt).toLocaleDateString()
                    : "N/A"}
                </p>
              </div>
            </div>
            <Button variant="outline" className="w-full">
              <ExternalLink className="h-4 w-4 mr-2" />
              Manage Account in Clerk
            </Button>
          </CardContent>
        </Card>

        {/* Appearance Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Palette className="h-5 w-5" />
              Appearance
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">Theme</p>
                <p className="text-sm text-muted-foreground">
                  Current theme: {themeName}
                </p>
              </div>
              <ThemeToggle />
            </div>
            <Separator />
            <div className="space-y-2">
              <p className="text-sm font-medium">Theme Information</p>
              <div className="space-y-1 text-sm text-muted-foreground">
                <p>• Light theme for daytime use</p>
                <p>• Dark theme for low-light environments</p>
                <p>• System theme follows your OS preference</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Notification Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Bell className="h-5 w-5" />
              Notifications
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Task Reminders</p>
                  <p className="text-sm text-muted-foreground">
                    Get notified about upcoming due dates
                  </p>
                </div>
                <Badge variant="outline">Coming Soon</Badge>
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Email Notifications</p>
                  <p className="text-sm text-muted-foreground">
                    Receive email updates about your tasks
                  </p>
                </div>
                <Badge variant="outline">Coming Soon</Badge>
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Browser Notifications</p>
                  <p className="text-sm text-muted-foreground">
                    Show desktop notifications
                  </p>
                </div>
                <Badge variant="outline">Coming Soon</Badge>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Privacy & Security */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5" />
              Privacy & Security
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-3">
              <div>
                <p className="font-medium">Data Encryption</p>
                <p className="text-sm text-muted-foreground">
                  Your data is encrypted both in transit and at rest
                </p>
              </div>
              <div>
                <p className="font-medium">Authentication</p>
                <p className="text-sm text-muted-foreground">
                  Powered by Clerk with industry-standard security
                </p>
              </div>
              <div>
                <p className="font-medium">Data Retention</p>
                <p className="text-sm text-muted-foreground">
                  Your data is stored securely and can be deleted upon request
                </p>
              </div>
            </div>
            <Separator />
            <Button variant="outline" className="w-full">
              <ExternalLink className="h-4 w-4 mr-2" />
              Privacy Policy
            </Button>
          </CardContent>
        </Card>

        {/* Data Management */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Database className="h-5 w-5" />
              Data Management
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-3">
              <div className="text-center p-4 border rounded-lg">
                <p className="text-2xl font-bold text-blue-600">0</p>
                <p className="text-sm font-medium">Total Tasks</p>
                <p className="text-xs text-muted-foreground">All time</p>
              </div>
              <div className="text-center p-4 border rounded-lg">
                <p className="text-2xl font-bold text-green-600">0</p>
                <p className="text-sm font-medium">Completed Tasks</p>
                <p className="text-xs text-muted-foreground">All time</p>
              </div>
              <div className="text-center p-4 border rounded-lg">
                <p className="text-2xl font-bold text-purple-600">0</p>
                <p className="text-sm font-medium">Categories</p>
                <p className="text-xs text-muted-foreground">Total created</p>
              </div>
            </div>
            <Separator />
            <div className="flex flex-col sm:flex-row gap-2">
              <Button variant="outline" className="flex-1">
                Export Data
              </Button>
              <Button variant="outline" className="flex-1">
                Import Data
              </Button>
              <Button variant="destructive" className="flex-1">
                Delete All Data
              </Button>
            </div>
            <p className="text-xs text-muted-foreground text-center">
              Data export/import and deletion features are coming soon. Contact
              support if you need immediate assistance.
            </p>
          </CardContent>
        </Card>
      </div>

      {/* App Information */}
      <Card>
        <CardHeader>
          <CardTitle>About Tasker</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <p className="font-medium">Version</p>
              <p className="text-sm text-muted-foreground">1.0.0</p>
            </div>
            <div>
              <p className="font-medium">Last Updated</p>
              <p className="text-sm text-muted-foreground">
                {new Date().toLocaleDateString()}
              </p>
            </div>
            <div>
              <p className="font-medium">Built With</p>
              <p className="text-sm text-muted-foreground">
                React, TypeScript, Tailwind CSS
              </p>
            </div>
            <div>
              <p className="font-medium">Support</p>
              <p className="text-sm text-muted-foreground">
                Contact us for help and feedback
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
