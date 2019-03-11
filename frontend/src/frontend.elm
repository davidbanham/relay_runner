module Main exposing (Model, init, main, subscriptions, update, view)

import Browser
import Browser.Navigation as Nav
import Debug exposing (log)
import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Json.Decode exposing (Decoder, bool, field, list)
import Msg exposing (Msg)
import RelayControl exposing (relayControl)
import String
import Time
import Types exposing (..)
import Url



-- MAIN


main : Program () Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        , onUrlChange = Msg.UrlChanged
        , onUrlRequest = Msg.LinkClicked
        }



-- MODEL


type alias Model =
    { relays : List Relay, err : String, currentTime : Time.Posix }


init : () -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
    ( Model [] "" (Time.millisToPosix 0)
    , Http.get
        { url = "/pins"
        , expect = Http.expectJson Msg.GotPins pinsDecoder
        }
    )



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Msg.GotPins result ->
            case result of
                Ok pins ->
                    ( { model
                        | relays = pins
                      }
                    , Cmd.none
                    )

                Err error ->
                    ( { model | err = Debug.toString error }, Cmd.none )

        Msg.Tick newTime ->
            ( { model | currentTime = newTime }
            , Cmd.none
            )

        Msg.LinkClicked _ ->
            ( model, Cmd.none )

        Msg.UrlChanged _ ->
            ( model, Cmd.none )

        Msg.OnUntil target value ->
            case String.toInt value of
                Nothing ->
                    case value == "" of
                        True ->
                            ( { model | err = "" }, Cmd.none )

                        False ->
                            ( { model | err = "Bad input: " ++ value }, Cmd.none )

                Just i ->
                    ( { model
                        | relays =
                            List.map
                                (\relay ->
                                    case relay.id == target.id of
                                        True ->
                                            { relay | setMinutes = i, onUntil = Time.millisToPosix (Time.posixToMillis model.currentTime + (i * 60 * 1000)) }

                                        False ->
                                            relay
                                )
                                model.relays
                        , err = ""
                      }
                    , Cmd.none
                    )

        Msg.ToggleRelay target state ->
            ( { model
                | relays =
                    List.map
                        (\relay ->
                            if relay.id == target.id then
                                { relay | on = state }

                            else
                                relay
                        )
                        model.relays
              }
            , Http.post
                { url =
                    "/pins/"
                        ++ String.fromInt target.id
                        ++ "/"
                        ++ (if state then
                                if target.setMinutes /= 0 then
                                    "on?length=" ++ String.fromInt target.setMinutes

                                else
                                    "on"

                            else
                                "off"
                           )
                , expect = Http.expectWhatever Msg.ToggledRelay
                , body = Http.emptyBody
                }
            )

        Msg.ToggledRelay _ ->
            ( model, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Time.every 1000 Msg.Tick



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "Relay Runner"
    , body =
        [ node "link"
            [ href "/css/main.css", rel "stylesheet" ]
            []
        , div [] [ text model.err ]
        , div [ class "wrapper" ]
            (List.map
                (relayControl model.currentTime)
                model.relays
            )
        ]
    }


pinsDecoder : Decoder (List Relay)
pinsDecoder =
    field "pins" (list Types.relayDecoder)
