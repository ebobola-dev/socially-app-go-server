ALTER TABLE privileges
ADD CONSTRAINT uq_privileges_order_index UNIQUE (order_index);
