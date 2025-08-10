import {
  useUpdateCategory,
  type TGetCategoriesResponse,
} from "@/api/hooks/use-category-query";
import { Badge } from "@/components/ui/badge";
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
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import { type ReactNode, useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import * as z from "zod";

const updateCategorySchema = z.object({
  name: z.string().min(1, "Name is required").max(100, "Name too long"),
  description: z.string().max(500, "Description too long").optional(),
  color: z.string().min(1, "Color is required"),
});

type UpdateCategoryForm = z.infer<typeof updateCategorySchema>;

interface CategoryEditFormProps {
  category: TGetCategoriesResponse["data"][number];
  children: ReactNode;
}

const predefinedColors = [
  "#ef4444", // red
  "#f97316", // orange
  "#eab308", // yellow
  "#22c55e", // green
  "#06b6d4", // cyan
  "#3b82f6", // blue
  "#8b5cf6", // violet
  "#ec4899", // pink
  "#6b7280", // gray
  "#84cc16", // lime
  "#f59e0b", // amber
  "#10b981", // emerald
  "#0ea5e9", // sky
  "#6366f1", // indigo
  "#a855f7", // purple
  "#e11d48", // rose
];

export function CategoryEditForm({
  category,
  children,
}: CategoryEditFormProps) {
  const [open, setOpen] = useState(false);
  const updateCategory = useUpdateCategory();

  const form = useForm<UpdateCategoryForm>({
    resolver: zodResolver(updateCategorySchema),
    defaultValues: {
      name: category.name,
      description: category.description || "",
      color: category.color,
    },
  });

  const selectedColor = form.watch("color");

  useEffect(() => {
    if (open) {
      form.reset({
        name: category.name,
        description: category.description || "",
        color: category.color,
      });
    }
  }, [open, category, form]);

  const onSubmit = async (data: UpdateCategoryForm) => {
    try {
      await updateCategory.mutateAsync({
        categoryId: category.id,
        body: {
          name: data.name,
          description: data.description || undefined,
          color: data.color,
        },
      });

      toast.success("Category updated successfully!");
      setOpen(false);
    } catch {
      toast.error("Failed to update category");
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Edit Category</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Enter category name..."
                      {...field}
                      disabled={updateCategory.isPending}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Enter category description..."
                      rows={3}
                      {...field}
                      disabled={updateCategory.isPending}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="color"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Color *</FormLabel>
                  <FormControl>
                    <div className="space-y-3">
                      <div className="flex items-center gap-3">
                        <div
                          className="w-8 h-8 rounded-full border-2 border-white shadow-sm"
                          style={{ backgroundColor: selectedColor }}
                        />
                        <Badge
                          variant="secondary"
                          className="flex items-center gap-2"
                        >
                          <div
                            className="w-2 h-2 rounded-full"
                            style={{ backgroundColor: selectedColor }}
                          />
                          {selectedColor}
                        </Badge>
                      </div>

                      <div className="grid grid-cols-8 gap-2">
                        {predefinedColors.map((color) => (
                          <button
                            key={color}
                            type="button"
                            className={`w-8 h-8 rounded-full border-2 hover:scale-110 transition-transform ${
                              selectedColor === color
                                ? "border-foreground shadow-md"
                                : "border-white shadow-sm"
                            }`}
                            style={{ backgroundColor: color }}
                            onClick={() => field.onChange(color)}
                            disabled={updateCategory.isPending}
                          />
                        ))}
                      </div>

                      <div className="flex items-center gap-2">
                        <Input
                          type="color"
                          value={selectedColor}
                          onChange={(e) => field.onChange(e.target.value)}
                          className="w-16 h-8 p-1 border rounded"
                          disabled={updateCategory.isPending}
                        />
                        <span className="text-sm text-muted-foreground">
                          Or pick a custom color
                        </span>
                      </div>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end gap-2 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                disabled={updateCategory.isPending}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={updateCategory.isPending}>
                {updateCategory.isPending ? "Updating..." : "Update Category"}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
