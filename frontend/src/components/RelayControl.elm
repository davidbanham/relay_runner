module RelayControl exposing (relayControl)

import Duration
import Html exposing (..)
import Html.Attributes exposing (checked, class, href, name, src, step, type_, value)
import Html.Events exposing (onCheck, onInput)
import Msg exposing (Msg)
import Time
import Types exposing (..)


noNeg : Float -> String
noNeg input =
    if input <= 0 then
        ""

    else
        String.fromInt (round input)


relayControl : Time.Posix -> Relay -> Html Msg
relayControl currentTime relay =
    div []
        [ h4 [] [ text ("Relay " ++ relay.name) ]
        , label
            [ class "switch"
            , onCheck (Msg.ToggleRelay relay)
            ]
            [ input [ checked relay.on, type_ "checkbox" ] []
            , span [ class "slider round" ] []
            ]
        , br [] []
        , input [ value (noNeg (Duration.inMinutes (Duration.from currentTime relay.onUntil))), onInput (Msg.OnUntil relay), class "timeinput", type_ "text", step "1", name "length" ] []
        ]
