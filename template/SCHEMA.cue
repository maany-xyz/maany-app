// CUE schema for template parameters

project_name:        string & !=""
module_path:         =~"^[a-zA-Z0-9._~\-/]+$"
binary_name:         =~"^[a-z0-9][a-z0-9-]{2,32}$"
bech32_main_prefix:  =~"^[a-z][a-z0-9]{2,16}$"
base_denom:          =~"^[a-z][a-z0-9/]{2,64}$"
display_denom:       =~"^[A-Z][A-Z0-9]{1,10}$"
denom_exponent:      >=0 & <=18
chain_id:            =~"^[a-z0-9-]{3,48}$"
min_gas_price:       =~"^[0-9]+(\.[0-9]+)?$"
home_dir_name:       =~"^\.[a-z0-9.-]{3,64}$"

