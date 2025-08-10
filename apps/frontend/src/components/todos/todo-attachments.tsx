import type { TGetTodoByIdResponse } from "@/api/hooks";
import {
  useDeleteTodoAttachment,
  useGetTodoAttachmentPresignedUrl,
} from "@/api/hooks/use-todo-attachment-query";
import { Button } from "@/components/ui/button";
import { formatFileSize } from "@/lib/utils";
import { Download, FileText, Trash2 } from "lucide-react";
import { toast } from "sonner";

interface TodoAttachmentsProps {
  todoId: string;
  attachments: TGetTodoByIdResponse["attachments"];
  disabled?: boolean;
}

export function TodoAttachments({
  todoId,
  attachments,
  disabled = false,
}: TodoAttachmentsProps) {
  const deleteAttachment = useDeleteTodoAttachment();
  const getPresignedUrl = useGetTodoAttachmentPresignedUrl();

  const handleDownload = async (
    attachment: TGetTodoByIdResponse["attachments"][number],
  ) => {
    try {
      const result = await getPresignedUrl.mutateAsync({
        todoId,
        attachmentId: attachment.id,
      });

      // Open the presigned URL in a new tab to trigger download
      window.open(result.url, "_blank");
    } catch (error) {
      console.error("Failed to download attachment:", error);
    }
  };

  const handleDelete = async (attachmentId: string, attachmentName: string) => {
    if (!confirm(`Are you sure you want to delete "${attachmentName}"?`)) {
      return;
    }

    try {
      await deleteAttachment.mutateAsync({
        todoId,
        attachmentId,
      });
      toast.success("Attachment deleted successfully");
    } catch (error) {
      console.error("Failed to delete attachment:", error);
    }
  };

  if (attachments.length === 0) {
    return null;
  }

  return (
    <div className="space-y-2">
      <h4 className="text-sm font-medium">
        Current Attachments ({attachments.length})
      </h4>
      <div className="space-y-2 max-h-48 overflow-y-auto">
        {attachments.map((attachment) => (
          <div
            key={attachment.id}
            className="flex items-center justify-between p-3 bg-muted/50 rounded-lg hover:bg-muted transition-colors"
          >
            <div className="flex items-center gap-3 min-w-0 flex-1">
              <FileText className="h-5 w-5 text-muted-foreground flex-shrink-0" />
              <div className="min-w-0 flex-1">
                <p
                  className="text-sm font-medium truncate"
                  title={attachment.name}
                >
                  {attachment.name}
                </p>
                <p className="text-xs text-muted-foreground">
                  {attachment.fileSize
                    ? formatFileSize(attachment.fileSize)
                    : "Unknown size"}
                  {attachment.mimeType &&
                    ` • ${attachment.mimeType.split("/")[1].toUpperCase()}`}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-1">
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => handleDownload(attachment)}
                disabled={disabled || getPresignedUrl.isPending}
                title="Download attachment"
              >
                <Download className="h-4 w-4" />
              </Button>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => handleDelete(attachment.id, attachment.name)}
                disabled={disabled || deleteAttachment.isPending}
                className="text-destructive hover:text-destructive"
                title="Delete attachment"
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
