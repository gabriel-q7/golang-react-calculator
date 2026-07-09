import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { OPERATIONS } from "../config";
import { OperationCard } from "./OperationCard";

export function CalculatorPage() {
  return (
    <div className="mx-auto w-full max-w-xl">
      <h1 className="mb-4 text-2xl font-bold tracking-tight">Calculator</h1>
      <Tabs defaultValue={OPERATIONS[0].id}>
        <TabsList>
          {OPERATIONS.map((op) => (
            <TabsTrigger key={op.id} value={op.id}>
              {op.symbol}
              <span className="ml-1 hidden sm:inline">{op.label}</span>
            </TabsTrigger>
          ))}
        </TabsList>
        {OPERATIONS.map((op) => (
          <TabsContent key={op.id} value={op.id}>
            <OperationCard config={op} />
          </TabsContent>
        ))}
      </Tabs>
    </div>
  );
}
