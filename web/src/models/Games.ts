interface Games{
    id : String
    quiz_id : String
    created_at : Date
    updated_at : Date
    code : String
    player_limit : BigInt
    start_time : Date
    finish_time : Date
}

export default Games