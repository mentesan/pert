CREATE TABLE "contact" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "phone" varchar(255) NOT NULL,
  "customer_id" int4 NOT NULL,
  CONSTRAINT "_copy_5" PRIMARY KEY ("id")
);

CREATE TABLE "customer" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "address" varchar(255),
  "site" varchar(255),
  "description" varchar(255),
  "info" varchar(255),
  CONSTRAINT "_copy_4" PRIMARY KEY ("id")
);

CREATE TABLE "framework" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" varchar(255),
  CONSTRAINT "_copy_8" PRIMARY KEY ("id")
);

CREATE TABLE "framework_phase" (
  "framework" int4 NOT NULL,
  "phase" int4 NOT NULL,
  "sequence" int4 NOT NULL,
  CONSTRAINT "_copy_7" PRIMARY KEY ("framework", "phase")
);

CREATE TABLE "image" (
  "id" int4 NOT NULL,
  "description" varchar(255),
  "path" varchar(255) NOT NULL,
  CONSTRAINT "_copy_10" PRIMARY KEY ("id")
);

CREATE TABLE "phase" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" varchar(255),
  "info" varchar(255),
  CONSTRAINT "_copy_6" PRIMARY KEY ("id")
);

CREATE TABLE "phase_tool" (
  "phase" int4 NOT NULL,
  "tool" int4 NOT NULL,
  CONSTRAINT "_copy_1" PRIMARY KEY ("tool", "phase")
);

CREATE TABLE "project" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "customer_id" int4 NOT NULL,
  "framework_id" int4 NOT NULL,
  "type" int2 NOT NULL,
  CONSTRAINT "_copy_17" PRIMARY KEY ("id")
);

CREATE TABLE "project_target" (
  "project" int4 NOT NULL,
  "target" int4 NOT NULL,
  CONSTRAINT "_copy_9" PRIMARY KEY ("project", "target")
);

CREATE TABLE "project_user" (
  "project" int4 NOT NULL,
  "user" int4 NOT NULL,
  CONSTRAINT "_copy_16" PRIMARY KEY ("user", "project")
);

CREATE TABLE "project_vuln" (
  "project" int4 NOT NULL,
  "vulnerability" int4 NOT NULL,
  CONSTRAINT "_copy_14" PRIMARY KEY ("project", "vulnerability")
);

CREATE TABLE "ProjectContacts" (
  "Project" int4 NOT NULL,
  "Contact" int4 NOT NULL,
  PRIMARY KEY ("Project", "Contact")
);

CREATE TABLE "risk" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "severity" int2 NOT NULL,
  "color" varchar(255) NOT NULL,
  "description" varchar(255),
  CONSTRAINT "_copy_12" PRIMARY KEY ("id")
);

CREATE TABLE "target" (
  "id" int4 NOT NULL,
  "name" varchar(255),
  "url" varchar(255),
  "ip" varchar(255) NOT NULL,
  "description" varchar(255),
  "network" varchar(255),
  "owner" varchar(255),
  "country" varchar(255),
  "state" varchar(255),
  "city" varchar(255),
  "traceroute" varchar(255),
  CONSTRAINT "_copy_3" PRIMARY KEY ("id")
);

CREATE TABLE "tool" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" varchar(255),
  "synopsis" varchar(255),
  "url" varchar(255),
  CONSTRAINT "_copy_2" PRIMARY KEY ("id")
);

CREATE TABLE "user" (
  "id" int4 NOT NULL,
  "category" varchar(255) NOT NULL,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "photo" varchar(255),
  "info" varchar(255) NOT NULL,
  CONSTRAINT "_copy_15" PRIMARY KEY ("id")
);

CREATE TABLE "vuln_image" (
  "vuln_id" int4 NOT NULL,
  "img_id" int4 NOT NULL,
  CONSTRAINT "_copy_11" PRIMARY KEY ("vuln_id", "img_id")
);

CREATE TABLE "vulnerability" (
  "id" int4 NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" varchar(255),
  "status" int2 NOT NULL,
  "poc_desc" varchar(255),
  "risk_id" int4 NOT NULL,
  CONSTRAINT "_copy_13" PRIMARY KEY ("id")
);

ALTER TABLE "contact" ADD CONSTRAINT "customer_copy_1" FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");
ALTER TABLE "framework_phase" ADD CONSTRAINT "framework_copy_1" FOREIGN KEY ("framework") REFERENCES "framework" ("id");
ALTER TABLE "framework_phase" ADD CONSTRAINT "phase" FOREIGN KEY ("phase") REFERENCES "phase" ("id");
ALTER TABLE "phase_tool" ADD CONSTRAINT "phase_copy_1" FOREIGN KEY ("phase") REFERENCES "phase" ("id");
ALTER TABLE "phase_tool" ADD CONSTRAINT "tool" FOREIGN KEY ("tool") REFERENCES "tool" ("id");
ALTER TABLE "project" ADD CONSTRAINT "framework" FOREIGN KEY ("framework_id") REFERENCES "framework" ("id");
ALTER TABLE "project" ADD CONSTRAINT "customer" FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");
ALTER TABLE "project_target" ADD CONSTRAINT "project_copy_1" FOREIGN KEY ("project") REFERENCES "project" ("id");
ALTER TABLE "project_target" ADD CONSTRAINT "target" FOREIGN KEY ("target") REFERENCES "target" ("id");
ALTER TABLE "project_user" ADD CONSTRAINT "project" FOREIGN KEY ("project") REFERENCES "project" ("id");
ALTER TABLE "project_user" ADD CONSTRAINT "user" FOREIGN KEY ("user") REFERENCES "user" ("id");
ALTER TABLE "project_vuln" ADD CONSTRAINT "project_copy_2" FOREIGN KEY ("project") REFERENCES "project" ("id");
ALTER TABLE "project_vuln" ADD CONSTRAINT "vulnerability" FOREIGN KEY ("vulnerability") REFERENCES "vulnerability" ("id");
ALTER TABLE "ProjectContacts" ADD CONSTRAINT "Project" FOREIGN KEY ("Project") REFERENCES "project" ("id");
ALTER TABLE "ProjectContacts" ADD CONSTRAINT "Contact" FOREIGN KEY ("Contact") REFERENCES "contact" ("id");
ALTER TABLE "vuln_image" ADD CONSTRAINT "vuln" FOREIGN KEY ("vuln_id") REFERENCES "vulnerability" ("id");
ALTER TABLE "vuln_image" ADD CONSTRAINT "img" FOREIGN KEY ("img_id") REFERENCES "image" ("id");
ALTER TABLE "vulnerability" ADD CONSTRAINT "risk" FOREIGN KEY ("risk_id") REFERENCES "risk" ("id");

