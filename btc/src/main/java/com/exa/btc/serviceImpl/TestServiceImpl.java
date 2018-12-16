package com.exa.btc.serviceImpl;

import com.btc.b1.service.TestService;
import org.springframework.stereotype.Service;

@Service
public class TestServiceImpl implements TestService {


    @Override
    public String hello() {
        return "bitcoinf   low";
    }
}
