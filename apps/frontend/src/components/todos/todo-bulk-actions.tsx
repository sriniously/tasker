import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useBulkUpdateTodos } from "@/api/hooks/use-todo-query";
import { useGetAllCategories } from "@/api/hooks/use-category-query";
import {
  MoreHorizontal,
  CheckCircle,
  Clock,
  Archive,
  Folder,
  Flag,
  Trash2,
} from "lucide-react";
import { toast } from "sonner";

interface TodoBulkActionsProps {
  selectedTodos: string[];
  onSelectionChange: (todos: string[]) => void;
}

export function TodoBulkActions({ selectedTodos, onSelectionChange }: TodoBulkActionsProps) {
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showStatusDialog, setShowStatusDialog] = useState(false);
  const [showPriorityDialog, setShowPriorityDialog] = useState(false);
  const [showCategoryDialog, setShowCategoryDialog] = useState(false);
  
  const [selectedStatus, setSelectedStatus] = useState<string>("");
  const [selectedPriority, setSelectedPriority] = useState<string>("");
  const [selectedCategory, setSelectedCategory] = useState<string>();

  const bulkUpdate = useBulkUpdateTodos();
  const { data: categories } = useGetAllCategories({
    query: { page: 1, limit: 100 },
  });

  const handleBulkStatusUpdate = async () => {
    if (!selectedStatus) return;
    
    try {
      await bulkUpdate.mutateAsync({
        body: {
          todoIds: selectedTodos,
          status: selectedStatus as "draft" | "active" | "completed" | "archived",
        },
      });
      
      toast.success(`${selectedTodos.length} tasks updated successfully!`);
      onSelectionChange([]);
      setShowStatusDialog(false);
      setSelectedStatus("");
    } catch {
      toast.error("Failed to update tasks");
    }
  };

  const handleBulkPriorityUpdate = async () => {
    if (!selectedPriority) return;
    
    try {
      await bulkUpdate.mutateAsync({
        body: {
          todoIds: selectedTodos,
          priority: selectedPriority as "low" | "medium" | "high",
        },
      });
      
      toast.success(`${selectedTodos.length} tasks updated successfully!`);
      onSelectionChange([]);
      setShowPriorityDialog(false);
      setSelectedPriority("");
    } catch {
      toast.error("Failed to update tasks");
    }
  };

  const handleBulkCategoryUpdate = async () => {
    try {
      await bulkUpdate.mutateAsync({
        body: {
          todoIds: selectedTodos,
          categoryId: selectedCategory || undefined,
        },
      });
      
      toast.success(`${selectedTodos.length} tasks updated successfully!`);
      onSelectionChange([]);
      setShowCategoryDialog(false);
      setSelectedCategory(undefined);
    } catch {
      toast.error("Failed to update tasks");
    }
  };

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="sm">
            <MoreHorizontal className="h-4 w-4 mr-2" />
            Bulk Actions
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-48">
          <DropdownMenuItem onClick={() => setShowStatusDialog(true)}>
            <CheckCircle className="h-4 w-4 mr-2" />
            Change Status
          </DropdownMenuItem>
          
          <DropdownMenuItem onClick={() => setShowPriorityDialog(true)}>
            <Flag className="h-4 w-4 mr-2" />
            Change Priority
          </DropdownMenuItem>
          
          <DropdownMenuItem onClick={() => setShowCategoryDialog(true)}>
            <Folder className="h-4 w-4 mr-2" />
            Change Category
          </DropdownMenuItem>
          
          <DropdownMenuSeparator />
          
          <DropdownMenuItem
            onClick={() => setShowDeleteDialog(true)}
            className="text-destructive focus:text-destructive"
          >
            <Trash2 className="h-4 w-4 mr-2" />
            Delete All
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Status Update Dialog */}
      <AlertDialog open={showStatusDialog} onOpenChange={setShowStatusDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Update Status</AlertDialogTitle>
            <AlertDialogDescription>
              Change the status for {selectedTodos.length} selected tasks.
            </AlertDialogDescription>
          </AlertDialogHeader>
          
          <div className="py-4">
            <Select value={selectedStatus} onValueChange={setSelectedStatus}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select new status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="draft">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-gray-500" />
                    Draft
                  </div>
                </SelectItem>
                <SelectItem value="active">
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-blue-500" />
                    Active
                  </div>
                </SelectItem>
                <SelectItem value="completed">
                  <div className="flex items-center gap-2">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    Completed
                  </div>
                </SelectItem>
                <SelectItem value="archived">
                  <div className="flex items-center gap-2">
                    <Archive className="h-4 w-4 text-gray-500" />
                    Archived
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setSelectedStatus("")}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleBulkStatusUpdate}
              disabled={!selectedStatus || bulkUpdate.isPending}
            >
              {bulkUpdate.isPending ? "Updating..." : "Update Status"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Priority Update Dialog */}
      <AlertDialog open={showPriorityDialog} onOpenChange={setShowPriorityDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Update Priority</AlertDialogTitle>
            <AlertDialogDescription>
              Change the priority for {selectedTodos.length} selected tasks.
            </AlertDialogDescription>
          </AlertDialogHeader>
          
          <div className="py-4">
            <Select value={selectedPriority} onValueChange={setSelectedPriority}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select new priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="low">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-green-500" />
                    Low
                  </div>
                </SelectItem>
                <SelectItem value="medium">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-yellow-500" />
                    Medium
                  </div>
                </SelectItem>
                <SelectItem value="high">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-red-500" />
                    High
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setSelectedPriority("")}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleBulkPriorityUpdate}
              disabled={!selectedPriority || bulkUpdate.isPending}
            >
              {bulkUpdate.isPending ? "Updating..." : "Update Priority"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Category Update Dialog */}
      <AlertDialog open={showCategoryDialog} onOpenChange={setShowCategoryDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Update Category</AlertDialogTitle>
            <AlertDialogDescription>
              Change the category for {selectedTodos.length} selected tasks.
            </AlertDialogDescription>
          </AlertDialogHeader>
          
          <div className="py-4">
            <Select value={selectedCategory} onValueChange={setSelectedCategory}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select new category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">No Category</SelectItem>
                {categories?.data?.map((category) => (
                  <SelectItem key={category.id} value={category.id}>
                    <div className="flex items-center gap-2">
                      <div
                        className="w-2 h-2 rounded-full"
                        style={{ backgroundColor: category.color }}
                      />
                      {category.name}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setSelectedCategory(undefined)}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleBulkCategoryUpdate}
              disabled={bulkUpdate.isPending}
            >
              {bulkUpdate.isPending ? "Updating..." : "Update Category"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Tasks</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete {selectedTodos.length} selected tasks? 
              This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                toast.error("Bulk delete feature coming soon!");
                setShowDeleteDialog(false);
              }}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete Tasks
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}