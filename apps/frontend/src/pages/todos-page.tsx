import { useGetAllCategories } from "@/api/hooks/use-category-query";
import { useDebounce } from "@/api/hooks/use-debounce";
import {
  useGetAllTodos,
  type TGetTodosQuery,
} from "@/api/hooks/use-todo-query";
import { TodoCard } from "@/components/todos/todo-card";
import { TodoCreateForm } from "@/components/todos/todo-create-form";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Plus, Search, CheckSquare } from "lucide-react";
import { useState } from "react";

type TodoStatus = "draft" | "active" | "completed" | "archived";
type TodoPriority = "low" | "medium" | "high" | "all";

export function TodosPage() {
  const [activeTab, setActiveTab] = useState<TodoStatus | "all">("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState<string>("all");
  const [selectedPriority, setSelectedPriority] = useState<TodoPriority>("all");
  const [sortBy, setSortBy] = useState<TGetTodosQuery["sort"]>("updated_at");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("desc");
  const [selectedTodos, setSelectedTodos] = useState<string[]>([]);
  const [page, setPage] = useState(1);

  const debouncedSearch = useDebounce(searchQuery, 300);

  const { data: todos, isLoading } = useGetAllTodos({
    query: {
      page,
      limit: 20,
      search: debouncedSearch || undefined,
      status: activeTab === "all" ? undefined : activeTab,
      categoryId: selectedCategory === "all" ? undefined : selectedCategory,
      priority: selectedPriority === "all" ? undefined : selectedPriority,
      sort: sortBy,
      order: sortOrder,
    },
  });

  const { data: categories } = useGetAllCategories({
    query: { page: 1, limit: 100 },
  });

  const handleSelectTodo = (todoId: string, checked: boolean) => {
    if (checked) {
      setSelectedTodos((prev) => [...prev, todoId]);
    } else {
      setSelectedTodos((prev) => prev.filter((id) => id !== todoId));
    }
  };

  const handleSelectAll = (checked: boolean) => {
    if (checked && todos?.data) {
      setSelectedTodos(todos.data.map((todo) => todo.id));
    } else {
      setSelectedTodos([]);
    }
  };

  const clearFilters = () => {
    setSearchQuery("");
    setSelectedCategory("all");
    setSelectedPriority("all");
    setSortBy("updated_at");
    setSortOrder("desc");
    setPage(1);
  };

  const tabCounts = {
    all: todos?.total || 0,
    draft: 0, // You can get these from stats if needed
    active: 0,
    completed: 0,
    archived: 0,
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Tasks</h1>
          <p className="text-muted-foreground">
            Manage and organize your tasks
          </p>
        </div>
        <TodoCreateForm>
          <Button>
            <Plus className="h-4 w-4 mr-2" />
            New Task
          </Button>
        </TodoCreateForm>
      </div>

      {/* Filters and Search */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Filters & Search</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search tasks..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
            <Select
              value={selectedCategory}
              onValueChange={setSelectedCategory}
            >
              <SelectTrigger className="w-full sm:w-48">
                <SelectValue placeholder="Category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                {categories?.data?.map((category) => (
                  <SelectItem key={category.id} value={category.id}>
                    <div className="flex items-center gap-2">
                      <div
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: category.color }}
                      />
                      {category.name}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select
              value={selectedPriority}
              onValueChange={(value) =>
                setSelectedPriority(value as TodoPriority)
              }
            >
              <SelectTrigger className="w-full sm:w-32">
                <SelectValue placeholder="Priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All</SelectItem>
                <SelectItem value="high">High</SelectItem>
                <SelectItem value="medium">Medium</SelectItem>
                <SelectItem value="low">Low</SelectItem>
              </SelectContent>
            </Select>
            <Select
              value={`${sortBy}-${sortOrder}`}
              onValueChange={(value) => {
                const [sort, order] = value.split("-");
                setSortBy(sort as TGetTodosQuery["sort"]);
                setSortOrder(order as "asc" | "desc");
              }}
            >
              <SelectTrigger className="w-full sm:w-40">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="updated_at-desc">Latest</SelectItem>
                <SelectItem value="updated_at-asc">Oldest</SelectItem>
                <SelectItem value="title-asc">Title A-Z</SelectItem>
                <SelectItem value="title-desc">Title Z-A</SelectItem>
                <SelectItem value="priority-desc">High Priority</SelectItem>
                <SelectItem value="due_date-asc">Due Date</SelectItem>
              </SelectContent>
            </Select>
            <Button variant="outline" onClick={clearFilters}>
              Clear
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Status Tabs */}
      <Tabs
        value={activeTab}
        onValueChange={(value) => setActiveTab(value as TodoStatus | "all")}
      >
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="all">All ({tabCounts.all})</TabsTrigger>
          <TabsTrigger value="draft">Draft ({tabCounts.draft})</TabsTrigger>
          <TabsTrigger value="active">Active ({tabCounts.active})</TabsTrigger>
          <TabsTrigger value="completed">
            Completed ({tabCounts.completed})
          </TabsTrigger>
          <TabsTrigger value="archived">
            Archived ({tabCounts.archived})
          </TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          {/* Bulk Actions */}
          {selectedTodos.length > 0 && (
            <Card>
              <CardContent className="py-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Checkbox
                      checked={selectedTodos.length === todos?.data?.length}
                      onCheckedChange={handleSelectAll}
                    />
                    <span className="text-sm font-medium">
                      {selectedTodos.length} selected
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Todo List */}
          <div className="space-y-3">
            {isLoading ? (
              <div className="space-y-3">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Card key={i}>
                    <CardContent className="p-4">
                      <div className="space-y-2">
                        <Skeleton className="h-4 w-full" />
                        <Skeleton className="h-3 w-2/3" />
                        <div className="flex gap-2">
                          <Skeleton className="h-5 w-16" />
                          <Skeleton className="h-5 w-16" />
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            ) : todos?.data?.length ? (
              <>
                {/* Select All Option */}
                {todos.data.length > 1 && (
                  <Card>
                    <CardContent className="py-3">
                      <div className="flex items-center gap-2">
                        <Checkbox
                          checked={selectedTodos.length === todos.data.length}
                          onCheckedChange={handleSelectAll}
                        />
                        <span className="text-sm text-muted-foreground">
                          Select all {todos.data.length} tasks
                        </span>
                      </div>
                    </CardContent>
                  </Card>
                )}

                {todos.data.map((todo) => (
                  <div key={todo.id} className="flex items-start gap-3">
                    <div className="pt-4">
                      <Checkbox
                        checked={selectedTodos.includes(todo.id)}
                        onCheckedChange={(checked) =>
                          handleSelectTodo(todo.id, checked as boolean)
                        }
                      />
                    </div>
                    <div className="flex-1">
                      <TodoCard todo={todo} />
                    </div>
                  </div>
                ))}

                {/* Pagination */}
                {todos.totalPages > 1 && (
                  <div className="flex items-center justify-between">
                    <p className="text-sm text-muted-foreground">
                      Showing {(page - 1) * 20 + 1} to{" "}
                      {Math.min(page * 20, todos.total)} of {todos.total} tasks
                    </p>
                    <div className="flex gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        disabled={page <= 1}
                        onClick={() => setPage(page - 1)}
                      >
                        Previous
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        disabled={page >= todos.totalPages}
                        onClick={() => setPage(page + 1)}
                      >
                        Next
                      </Button>
                    </div>
                  </div>
                )}
              </>
            ) : (
              <Card>
                <CardContent className="py-12 text-center">
                  <CheckSquare className="h-12 w-12 mx-auto mb-4 text-muted-foreground opacity-50" />
                  <h3 className="text-lg font-medium mb-2">No tasks found</h3>
                  <p className="text-muted-foreground mb-4">
                    {searchQuery || selectedCategory || selectedPriority
                      ? "Try adjusting your filters or search terms."
                      : "Create your first task to get started!"}
                  </p>
                  {!searchQuery && !selectedCategory && !selectedPriority && (
                    <TodoCreateForm>
                      <Button>
                        <Plus className="h-4 w-4 mr-2" />
                        Create Task
                      </Button>
                    </TodoCreateForm>
                  )}
                </CardContent>
              </Card>
            )}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
