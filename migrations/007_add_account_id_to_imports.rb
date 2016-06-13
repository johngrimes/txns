Sequel.migration do
  up do
    alter_table :imports do
      add_foreign_key :account_id, :accounts
    end
  end
end
