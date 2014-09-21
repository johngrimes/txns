migration 'create txns' do
  database.create_table :txns do
    primary_key :id
    column :hash, 'varchar(128)', :null => false
    date :date, :null => false
    column :description, 'varchar(255)'
    integer :debit_cents
    integer :credit_cents
    integer :balance_cents, :null => false

    index :hash, :unique => true
  end
end
