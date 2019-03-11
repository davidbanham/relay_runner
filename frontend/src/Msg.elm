module Msg exposing (Msg(..))

import Browser
import Http
import Time
import Types exposing (..)
import Url


type Msg
    = ToggleRelay Relay Bool
    | OnUntil Relay String
    | LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | Tick Time.Posix
    | GotPins (Result Http.Error (List Relay))
    | ToggledRelay (Result Http.Error ())
