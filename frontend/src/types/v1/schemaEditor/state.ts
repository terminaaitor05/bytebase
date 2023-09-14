import { ComposedDatabase, ComposedProject } from "..";
import { Engine } from "../../proto/v1/common";
import { SchemaDesign } from "../../proto/v1/schema_design_service";
import { Schema } from "./atomType";

export enum SchemaEditorTabType {
  TabForDatabase = "database",
  TabForTable = "table",
}

// Tab context for editing database.
export interface DatabaseTabContext {
  id: string;
  type: SchemaEditorTabType.TabForDatabase;
  // Parent could be a database.
  parentName: string;
  selectedSchemaId?: string;
}

// Tab context for editing table.
export interface TableTabContext {
  id: string;
  type: SchemaEditorTabType.TabForTable;
  // Parent could be a database or a branch.
  parentName: string;
  schemaId: string;
  tableId: string;
}

export type TabContext = DatabaseTabContext | TableTabContext;

export interface DatabaseSchema {
  database: ComposedDatabase;
  schemaList: Schema[];
  originSchemaList: Schema[];
}

export interface BranchSchema {
  branch: SchemaDesign;
  schemaList: Schema[];
  originSchemaList: Schema[];
}

export interface SchemaEditorV1State {
  engine: Engine;
  project: ComposedProject;
  readonly: boolean;
  resourceType: "database" | "branch";
  resourceMap: {
    database: Map<string, DatabaseSchema>;
    branch: Map<string, BranchSchema>;
  };
  tabState: {
    tabMap: Map<string, TabContext>;
    currentTabId?: string;
  };
}