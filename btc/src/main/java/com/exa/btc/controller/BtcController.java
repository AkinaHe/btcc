package com.exa.btc.controller;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
public class BtcController {
    @RequestMapping("/btc")
    @ResponseBody
    public String he(){
        return  "btc";
    }

}
