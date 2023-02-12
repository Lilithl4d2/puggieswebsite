package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../queue"
)


func HttpReqLeaveQueue(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (!pPlayer.IsInQueue) {
				mapResponse["error"] = "You are not in queue";
			} else {
				//Leave queue
				if (pPlayer.DuoWith != "") {
					pPlayer2, bFound := players.MapPlayers[pPlayer.DuoWith];
					if (bFound && pPlayer2 != nil) {
						queue.Leave(pPlayer2, false);
					}
				}
				queue.Leave(pPlayer, false);
				mapResponse["success"] = true;
			}
			players.MuPlayers.Unlock();
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Please authorize first";
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
