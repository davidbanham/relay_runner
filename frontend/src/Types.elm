module Types exposing (Relay, relayDecoder)

import Iso8601
import Json.Decode exposing (Decoder, andThen, bool, fail, field, int, list, map4, string, succeed)
import Json.Decode.Pipeline as D
import Time


type alias Relay =
    { name : String
    , on : Bool
    , id : Int
    , onUntil : Time.Posix
    , setMinutes : Int
    }


relayDecoder : Decoder Relay
relayDecoder =
    succeed Relay
        |> D.required "name" string
        |> D.required "state" bool
        |> D.required "id" int
        |> D.required "stop_at" Iso8601.decoder
        |> D.hardcoded 0
