import { type ReactNode, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useUser } from "@clerk/clerk-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
import {
  useGetCommentsByTodoId,
  useAddComment,
  useUpdateComment,
  useDeleteComment,
} from "@/api/hooks/use-comment-query";
import {
  Send,
  MessageSquare,
  MoreHorizontal,
  Edit,
  Trash2,
} from "lucide-react";
import { toast } from "sonner";

const commentSchema = z.object({
  content: z.string().min(1, "Comment cannot be empty").max(1000, "Comment too long"),
});

type CommentForm = z.infer<typeof commentSchema>;

interface TodoCommentsDialogProps {
  todoId: string;
  children: ReactNode;
}

interface Comment {
  id: string;
  content: string;
  userId: string;
  createdAt: string;
  updatedAt: string;
}

export function TodoCommentsDialog({ todoId, children }: TodoCommentsDialogProps) {
  const [open, setOpen] = useState(false);
  const [editingComment, setEditingComment] = useState<string | null>(null);
  const [deleteCommentId, setDeleteCommentId] = useState<string | null>(null);
  
  const { user } = useUser();
  const { data: comments, isLoading } = useGetCommentsByTodoId({ todoId });
  const addComment = useAddComment();
  const updateComment = useUpdateComment();
  const deleteComment = useDeleteComment();

  const form = useForm<CommentForm>({
    resolver: zodResolver(commentSchema),
    defaultValues: {
      content: "",
    },
  });

  const editForm = useForm<CommentForm>({
    resolver: zodResolver(commentSchema),
    defaultValues: {
      content: "",
    },
  });

  const onSubmit = async (data: CommentForm) => {
    try {
      await addComment.mutateAsync({
        body: {
          todoId,
          content: data.content,
        },
      });
      
      toast.success("Comment added successfully!");
      form.reset();
    } catch {
      toast.error("Failed to add comment");
    }
  };

  const onEditSubmit = async (data: CommentForm) => {
    if (!editingComment) return;
    
    try {
      await updateComment.mutateAsync({
        commentId: editingComment,
        body: {
          content: data.content,
        },
      });
      
      toast.success("Comment updated successfully!");
      setEditingComment(null);
      editForm.reset();
    } catch {
      toast.error("Failed to update comment");
    }
  };

  const handleDeleteComment = async (commentId: string) => {
    try {
      await deleteComment.mutateAsync({ commentId });
      toast.success("Comment deleted successfully!");
      setDeleteCommentId(null);
    } catch {
      toast.error("Failed to delete comment");
    }
  };

  const startEditing = (comment: Comment) => {
    setEditingComment(comment.id);
    editForm.setValue("content", comment.content);
  };

  const cancelEditing = () => {
    setEditingComment(null);
    editForm.reset();
  };

  return (
    <>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          {children}
        </DialogTrigger>
        <DialogContent className="sm:max-w-[600px] h-[80vh] flex flex-col">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <MessageSquare className="h-5 w-5" />
              Comments ({comments?.length || 0})
            </DialogTitle>
          </DialogHeader>
          
          <div className="flex-1 flex flex-col min-h-0">
            {/* Comments List */}
            <ScrollArea className="flex-1 pr-4">
              <div className="space-y-4">
                {isLoading ? (
                  <div className="space-y-4">
                    {Array.from({ length: 3 }).map((_, i) => (
                      <Card key={i}>
                        <CardContent className="p-4">
                          <div className="space-y-2">
                            <div className="flex items-center justify-between">
                              <Skeleton className="h-4 w-24" />
                              <Skeleton className="h-4 w-16" />
                            </div>
                            <Skeleton className="h-12 w-full" />
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                ) : comments?.length ? (
                  comments.map((comment) => (
                    <Card key={comment.id}>
                      <CardContent className="p-4">
                        <div className="space-y-3">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <Badge variant="secondary">You</Badge>
                              <span className="text-xs text-muted-foreground">
                                {comment.createdAt !== comment.updatedAt && "edited • "}
                                {new Date(comment.createdAt).toLocaleString()}
                              </span>
                            </div>
                            
                            {comment.userId === user?.id && (
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button variant="ghost" size="icon" className="h-6 w-6">
                                    <MoreHorizontal className="h-4 w-4" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem onClick={() => startEditing(comment)}>
                                    <Edit className="h-4 w-4 mr-2" />
                                    Edit
                                  </DropdownMenuItem>
                                  <DropdownMenuItem
                                    onClick={() => setDeleteCommentId(comment.id)}
                                    className="text-destructive focus:text-destructive"
                                  >
                                    <Trash2 className="h-4 w-4 mr-2" />
                                    Delete
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            )}
                          </div>
                          
                          {editingComment === comment.id ? (
                            <Form {...editForm}>
                              <form onSubmit={editForm.handleSubmit(onEditSubmit)} className="space-y-2">
                                <FormField
                                  control={editForm.control}
                                  name="content"
                                  render={({ field }) => (
                                    <FormItem>
                                      <FormControl>
                                        <Textarea
                                          placeholder="Edit your comment..."
                                          rows={3}
                                          {...field}
                                          disabled={updateComment.isPending}
                                        />
                                      </FormControl>
                                      <FormMessage />
                                    </FormItem>
                                  )}
                                />
                                <div className="flex justify-end gap-2">
                                  <Button
                                    type="button"
                                    variant="outline"
                                    size="sm"
                                    onClick={cancelEditing}
                                    disabled={updateComment.isPending}
                                  >
                                    Cancel
                                  </Button>
                                  <Button
                                    type="submit"
                                    size="sm"
                                    disabled={updateComment.isPending}
                                  >
                                    {updateComment.isPending ? "Updating..." : "Update"}
                                  </Button>
                                </div>
                              </form>
                            </Form>
                          ) : (
                            <p className="text-sm whitespace-pre-wrap">{comment.content}</p>
                          )}
                        </div>
                      </CardContent>
                    </Card>
                  ))
                ) : (
                  <div className="text-center py-8 text-muted-foreground">
                    <MessageSquare className="h-12 w-12 mx-auto mb-2 opacity-50" />
                    <p>No comments yet.</p>
                    <p className="text-xs">Be the first to add a comment!</p>
                  </div>
                )}
              </div>
            </ScrollArea>

            {/* Add Comment Form */}
            <div className="border-t pt-4 mt-4">
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                  <FormField
                    control={form.control}
                    name="content"
                    render={({ field }) => (
                      <FormItem>
                        <FormControl>
                          <Textarea
                            placeholder="Add a comment..."
                            rows={3}
                            {...field}
                            disabled={addComment.isPending}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <div className="flex justify-end">
                    <Button type="submit" disabled={addComment.isPending}>
                      <Send className="h-4 w-4 mr-2" />
                      {addComment.isPending ? "Adding..." : "Add Comment"}
                    </Button>
                  </div>
                </form>
              </Form>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      {/* Delete Comment Dialog */}
      <AlertDialog
        open={!!deleteCommentId}
        onOpenChange={(open) => !open && setDeleteCommentId(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Comment</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this comment? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => deleteCommentId && handleDeleteComment(deleteCommentId)}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={deleteComment.isPending}
            >
              {deleteComment.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}